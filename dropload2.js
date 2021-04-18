/**
 * dropload:www.zugaya.ml
 * Linbo Zugaya
 */

( function ( doc, undefined ) {
    'use strict';
    var win = window;
    var doc = document;
    function DropLoad(element, options){
        this.element = element;
        // 上方是否插入DOM
        this.upInsertDOM = false;
        // loading状态
        this.loading = false;
        // 是否锁定
        this.isLockUp = false;
        this.isLockDown = false;
        // 是否有数据
        this.isData = true;
        this._scrollTop = 0;
        this._threshold = 0;
        this.init(options);
    }
	
	function wscrollTop(){
		var t=document.documentElement.scrollTop;
		if(t)return t;
		else return document.body.scrollTop;
	}

    // 初始化
    DropLoad.prototype.init = function(options){
        var me = this;
		if(options.scrollArea===undefined)options.scrollArea = me.element;                                            // 滑动区域
		if(options.domUp===undefined)options.domUp = {                                                            // 上方DOM
                domClass   : 'dropload-up',
                domRefresh : '<div class="dropload-refresh">↓下拉刷新</div>',
                domUpdate  : '<div class="dropload-update">↑释放更新</div>',
                domLoad    : '<div class="dropload-load"><span class="loading"></span>加载中...</div>'
            };
        if(options.domDown===undefined)options.domDown = {                                                          // 下方DOM
                domClass   : 'dropload-down',
                domRefresh : '<div class="dropload-refresh">↑上拉加载更多</div>',
                domLoad    : '<div class="dropload-load"><span class="loading"></span>加载中...</div>',
                domNoData  : '<div class="dropload-noData">暂无数据</div>'
            };
        if(options.autoLoad===undefined)options.autoLoad = true;                                                     // 自动加载
        if(options.distance===undefined)options.distance = 50;                                                       // 拉动距离
        if(options.threshold===undefined)options.threshold = '';                                                      // 提前加载距离
        if(options.loadUpFn===undefined)options.loadUpFn = '';                                                       // 上方function
        if(options.loadDownFn===undefined)options.loadDownFn = '';                                                      // 下方function
		me.opts = options;

        // 如果加载下方，事先在下方插入DOM
        if(me.opts.loadDownFn != ''){
            var downnode=h('<div class="'+me.opts.domDown.domClass+'">'+me.opts.domDown.domRefresh+'</div>');
			me.element.appendChild(downnode);
            me.domDown = downnode;
        }

        // 计算提前加载距离
        if(!!me.domDown && me.opts.threshold === ''){
            // 默认滑到加载区2/3处时加载
            me._threshold = Math.floor(me.domDown.offsetHeight*1/3);
        }else{
            me._threshold = me.opts.threshold;
        }
		me._threshold=50;

        // 判断滚动区域
        if(me.opts.scrollArea == win){
            me.scrollArea = window;
            // 获取文档高度
            me._scrollContentHeight = document.body.scrollHeight;
            // 获取win显示区高度  —— 这里有坑
            me._scrollWindowHeight = doc.documentElement.clientHeight;
        }else{
            me.scrollArea = me.opts.scrollArea;
            me._scrollContentHeight = me.element.scrollHeight;
            me._scrollWindowHeight = me.element.offsetHeight;
        }
        fnAutoLoad(me);

        // 绑定触摸
        me.element.ontouchstart=function(e){
            if(!me.loading){
                fnTouches(e);
                fnTouchstart(e, me);
            }
        };
        me.element.ontouchmove=function(e){
            if(!me.loading){
                fnTouches(e, me);
                fnTouchmove(e, me);
            }
        };
        me.element.ontouchend=function(){
            if(!me.loading){
                fnTouchend(me);
            }
        };

        // 加载下方
        me.scrollArea.onscroll=function(){
			fnRecoverContentHeight(me);
			if(me.opts.scrollArea == win){
				me._scrollTop = wscrollTop();
			}else{
				me._scrollTop = me.scrollArea.scrollTop;
			}
            // 滚动页面触发加载数据
			console.log(me._scrollContentHeight+" "+me._threshold+"  "+me._scrollWindowHeight +"  "+ me._scrollTop+(me.opts.scrollArea == win)+" "+document.body.scrollTop);
            if(me.opts.loadDownFn != '' && !me.loading && !me.isLockDown && (me._scrollContentHeight - me._threshold) <= (me._scrollWindowHeight + me._scrollTop+10)){
                loadDown(me);
            }
        }
    };
	
	// 窗口调整
	DropLoad.prototype.startLoad=function(){
		var me=this;
		clearTimeout(me.timer);
		me.timer = setTimeout(function(){
			if(me.opts.scrollArea == win){
			// 重新获取win显示区高度
			me._scrollWindowHeight = win.innerHeight;
			}else{
				me._scrollWindowHeight = me.element.offsetHeight;
			}
			fnAutoLoad(me);
		},150);
	};

    // touches
    function fnTouches(e){
        if(!e.touches){
            e.touches = e.originalEvent.touches;
        }
    }

    // touchstart
    function fnTouchstart(e, me){
        me._startY = e.touches[0].pageY;
        // 记住触摸时的scrolltop值
		if(me.opts.scrollArea == win){
			me.touchScrollTop = wscrollTop();
		}else{
			me.touchScrollTop = me.scrollArea.scrollTop;
		}
    }

    // touchmove
    function fnTouchmove(e, me){
        me._curY = e.touches[0].pageY;
        me._moveY = me._curY - me._startY;

        if(me._moveY > 0){
            me.direction = 'down';
        }else if(me._moveY < 0){
            me.direction = 'up';
        }

        var _absMoveY = Math.abs(me._moveY);

        // 加载上方
        if(me.opts.loadUpFn != '' && me.touchScrollTop <= 0 && me.direction == 'down' && !me.isLockUp){
            e.preventDefault();

            me.domUp = c(me.opts.domUp.domClass);
            // 如果加载区没有DOM
            if(!me.upInsertDOM){
                me.element.innerHTML='<div class="'+me.opts.domUp.domClass+'"></div>'+me.element.innerHTML;
                me.upInsertDOM = true;
            }
            
            fnTransition(me.domUp,0);

            // 下拉
            if(_absMoveY <= me.opts.distance){
                me._offsetY = _absMoveY;
                // todo：move时会不断清空、增加dom，有可能影响性能，下同
                me.domUp.innerHTML=me.opts.domUp.domRefresh;
            // 指定距离 < 下拉距离 < 指定距离*2
            }else if(_absMoveY > me.opts.distance && _absMoveY <= me.opts.distance*2){
                me._offsetY = me.opts.distance+(_absMoveY-me.opts.distance)*0.5;
                me.domUp.innerHTML=me.opts.domUp.domUpdate;
            // 下拉距离 > 指定距离*2
            }else{
                me._offsetY = me.opts.distance+me.opts.distance*0.5+(_absMoveY-me.opts.distance*2)*0.2;
            }

            me.domUp.style.height=me._offsetY;
        }
    }

    // touchend
    function fnTouchend(me){
        var _absMoveY = Math.abs(me._moveY);
        if(me.opts.loadUpFn != '' && me.touchScrollTop <= 0 && me.direction == 'down' && !me.isLockUp){
            fnTransition(me.domUp,300);

            if(_absMoveY > me.opts.distance){
                me.domUp.style.height=me.domUp.children[0].offsetHeight;
                me.domUp.innerHTML=me.opts.domUp.domLoad;
                me.loading = true;
                me.opts.loadUpFn(me);
            }else{
                me.domUp.style.height=0;
				me.domUp.onwebkitTransitionEnd=
				me.domUp.onmozTransitionEnd=
				me.domUp.ontransitionend=function(){
                    me.upInsertDOM = false;
                    this.parentNode.removeChild(this);
                };
            }
            me._moveY = 0;
        }
    }

    // 如果文档高度不大于窗口高度，数据较少，自动加载下方数据
    function fnAutoLoad(me){
		fnRecoverContentHeight(me);
		console.log(me.opts.loadDownFn+"  "+me.opts.autoLoad);
        if(me.opts.loadDownFn != '' && me.opts.autoLoad){
			console.log(me._scrollContentHeight+" "+me._threshold+"  "+me._scrollWindowHeight);
            if((me._scrollContentHeight - me._threshold) <= me._scrollWindowHeight){
                loadDown(me);
            }
        }
    }

    // 重新获取文档高度
    function fnRecoverContentHeight(me){
        if(me.opts.scrollArea == win){
            me._scrollContentHeight = document.body.scrollHeight;
        }else{
			console.log("me.element.scrollHeight:",me.element.firstChild.children[1].offsetHeight)
            me._scrollContentHeight = me.element.firstChild.children[1].offsetHeight;
        }
    }

    // 加载下方
    function loadDown(me){
        me.direction = 'up';
        me.domDown.innerHTML=me.opts.domDown.domLoad;
        me.loading = true;
        me.opts.loadDownFn(me);
    }

    // 锁定
    DropLoad.prototype.lock = function(direction){
        var me = this;
        // 如果不指定方向
        if(direction === undefined){
            // 如果操作方向向上
            if(me.direction == 'up'){
                me.isLockDown = true;
            // 如果操作方向向下
            }else if(me.direction == 'down'){
                me.isLockUp = true;
            }else{
                me.isLockUp = true;
                me.isLockDown = true;
            }
        // 如果指定锁上方
        }else if(direction == 'up'){
            me.isLockUp = true;
        // 如果指定锁下方
        }else if(direction == 'down'){
            me.isLockDown = true;
            // 为了解决DEMO5中tab效果bug，因为滑动到下面，再滑上去点tab，direction=down，所以有bug
            me.direction = 'up';
        }
    };

    // 解锁
    DropLoad.prototype.unlock = function(){
        var me = this;
        // 简单粗暴解锁
        me.isLockUp = false;
        me.isLockDown = false;
        // 为了解决DEMO5中tab效果bug，因为滑动到下面，再滑上去点tab，direction=down，所以有bug
        me.direction = 'up';
    };

    // 无数据
    DropLoad.prototype.noData = function(flag){
        var me = this;
        if(flag === undefined || flag == true){
            me.isData = false;
        }else if(flag == false){
            me.isData = true;
        }
    };

    // 重置
    DropLoad.prototype.resetload = function(){
        var me = this;
        if(me.direction == 'down' && me.upInsertDOM){
            me.domUp.style.height=0;
			me.domUp.onwebkitTransitionEnd =
			me.domUp.onmozTransitionEnd =
			me.domUp.ontransitionend=function(){
                me.loading = false;
                me.upInsertDOM = false;
                this.parentNode.removeChild(this);
                fnRecoverContentHeight(me);
            };
        }else if(me.direction == 'up'){
            me.loading = false;
            // 如果有数据
            if(me.isData){
                // 加载区修改样式
                me.domDown.innerHTML=me.opts.domDown.domRefresh;
                fnRecoverContentHeight(me);
                fnAutoLoad(me);
            }else{
                // 如果没数据
                me.domDown.innerHTML=me.opts.domDown.domNoData;
            }
        }
    };

    // css过渡
    function fnTransition(dom,num){
        dom.css({
            '-webkit-transition':'all '+num+'ms',
            'transition':'all '+num+'ms'
        });
    }
	
	if ( typeof exports === 'object' ) {
		module.exports = DropLoad;
	} else if ( typeof define === 'function' && define.amd ) {
		define( function () {
			return DropLoad;
		});
	} else {
		win.DropLoad = DropLoad;
	}
}( document ) );