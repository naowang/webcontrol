package webcontrol

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"toolfunc"

	"net/http"

	"github.com/emirpasic/gods/maps/treemap"
)

const (
	Country_English = iota
	Country_China
)

const (
	T_FONTSIZE = iota
	T_MAKEHEADER
	T_TEXTALIGNMENT
	T_MAKELINK
	T_MAKELIST
	T_INSERTFILE
	T_UNDO
	T_REDO
	T_FONT
	T_FONTCOLOR
	T_BACKGROUNDCOLOR
	T_OK
	T_IMAGEFILE
	T_KEYWORDS
	T_UPLOAD
	T_LINK
	T_SEPARATORBYCOMMA
	T_FORMAT
	T_EXAMPLE
	T_MORE
	T_Menu_2
	T_PleaseSelectOneFile
)

var TranslateMap sync.Map
var IpCountryMap *treemap.Map
var SessionMap *sync.Map

func init() {
	TranslateMap.Store(T_FONTSIZE, []string{"font size", "字体大小"})
	TranslateMap.Store(T_MAKEHEADER, []string{"make header", "当标题"})
	TranslateMap.Store(T_TEXTALIGNMENT, []string{"text alignment", "文本对齐"})
	TranslateMap.Store(T_MAKELINK, []string{"make link", "作链接"})
	TranslateMap.Store(T_MAKELIST, []string{"make list", "作列表"})
	TranslateMap.Store(T_INSERTFILE, []string{"insert file", "插入文件"})
	TranslateMap.Store(T_UNDO, []string{"undo", "撤销"})
	TranslateMap.Store(T_REDO, []string{"redo", "重做"})
	TranslateMap.Store(T_FONT, []string{"font", "字体"})
	TranslateMap.Store(T_FONTCOLOR, []string{"font color", "字体颜色"})
	TranslateMap.Store(T_BACKGROUNDCOLOR, []string{"background color", "背景颜色"})
	TranslateMap.Store(T_OK, []string{"OK", "确定"})
	TranslateMap.Store(T_IMAGEFILE, []string{"image file", "图片文件"})
	TranslateMap.Store(T_KEYWORDS, []string{"keywords", "关键字"})
	TranslateMap.Store(T_UPLOAD, []string{"upload", "上传"})
	TranslateMap.Store(T_LINK, []string{"link", "链接"})
	TranslateMap.Store(T_SEPARATORBYCOMMA, []string{"separator by comma", "逗号分割"})
	TranslateMap.Store(T_FORMAT, []string{"format", "格式化"})
	TranslateMap.Store(T_EXAMPLE, []string{"example", "例子"})
	TranslateMap.Store(T_MORE, []string{"more", "更多"})
	TranslateMap.Store(T_Menu_2, []string{"menu 2", "菜单2"})
	TranslateMap.Store(T_PleaseSelectOneFile, []string{"Please Select 1 File", "请选择1个文件"})

	IpCountryMap = treemap.NewWithStringComparator()
	cnipsctt, _ := toolfunc.ReadFile(`M:\work\code\go\src\webcontrol\countryips\cn_ips.txt`)
	linels := bytes.Split(cnipsctt, []byte{'\n'})
	for i := 0; i < len(linels); i++ {
		linsegments := bytes.Split(linels[i], []byte{'\t'})
		if len(linsegments) == 4 {
			IpCountryMap.Put(string(linsegments[0]), 1)
		}
	}
}

func GetCountryId(r *http.Request) int {
	ip := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
	return IpToCountryId(ip)
}

func IpToCountryId(ip string) int {
	fndkey, fndval := IpCountryMap.Floor(ip)
	if fndkey != nil {
		return fndval.(int)
	}
	return 0
}

func GetParam(r *http.Request, name string) string {
	value2 := r.PostFormValue(name)
	if len(value2) > 0 {
		return value2
	}
	result := r.URL.Query().Get(name)
	return result
}

const (
	TOP_BAR    = 1
	RIGHT_BAR  = 2
	BOTTOM_BAR = 3
	LEFT_BAR   = 4
)

func T(textid, countryid int) string {
	val, bval := TranslateMap.Load(textid)
	if bval {
		if countryid < len(val.([]string)) {
			return val.([]string)[countryid]
		} else if len(val.([]string)) > 0 {
			return val.([]string)[0]
		}
	}
	return "unknow" + toolfunc.IntToStr(textid)
}

func CombineJs(html, title, stylesheet, selfjscode, jsfiletotal string, countryid int) []byte {
	var styletotal string = "html,body{width:100%;height:100%;padding:0px;margin:0px;}"
	if strings.Index(html, "setFontSize:before") != -1 {
		styletotal += `
@font-face            
{
font-family: 'fontawesome';
src: url('/font/fontawesome-webfont.woff') format('woff');
}
		`
	}
	styletotal += stylesheet
	jstotal := `function d(idstr){return document.getElementById(idstr);}
function c(csstr){return document.getElementsByClassName(csstr);}
function h(html){var div = document.createElement('div');div.innerHTML=html;return div.firstChild;}
function submit(URL, PARAMTERS,newWindow){var temp_form = document.createElement("form");if(newWindow){temp_form.target = '_blank';}temp_form.action = URL;temp_form.method = "post";temp_form.style.display = "none";for (var item in PARAMTERS) {var opt = document.createElement("textarea");opt.name = item;opt.value = PARAMTERS[item];temp_form.appendChild(opt);}document.body.appendChild(temp_form);temp_form.submit();}
        
!function(e,t){"object"==typeof exports&&"undefined"!=typeof module?module.exports=t():"function"==typeof define&&define.amd?define(t):(e=e||self,function(){var n=e.Cookies,r=e.Cookies=t();r.noConflict=function(){return e.Cookies=n,r}}())}(this,function(){"use strict";var e={read:function(e){return e.replace(/(%[\dA-F]{2})+/gi,decodeURIComponent)},write:function(e){return encodeURIComponent(e).replace(/%(2[346BF]|3[AC-F]|40|5[BDE]|60|7[BCD])/g,decodeURIComponent)}};function t(e){for(var t=1;t<arguments.length;t++){var n=arguments[t];for(var r in n)e[r]=n[r]}return e}return function n(r,o){function i(e,n,i){if("undefined"!=typeof document){"number"==typeof(i=t({},o,i)).expires&&(i.expires=new Date(Date.now()+864e5*i.expires)),i.expires&&(i.expires=i.expires.toUTCString()),n=r.write(n,e),e=encodeURIComponent(e).replace(/%(2[346B]|5E|60|7C)/g,decodeURIComponent).replace(/[()]/g,escape);var c="";for(var u in i)i[u]&&(c+="; "+u,!0!==i[u]&&(c+="="+i[u].split(";")[0]));return document.cookie=e+"="+n+c}}return Object.create({set:i,get:function(t){if("undefined"!=typeof document&&(!arguments.length||t)){for(var n=document.cookie?document.cookie.split("; "):[],o={},i=0;i<n.length;i++){var c=n[i].split("="),u=c.slice(1).join("=");'"'===u[0]&&(u=u.slice(1,-1));try{var f=e.read(c[0]);if(o[f]=r.read(u,f),t===f)break}catch(e){}}return t?o[t]:o}},remove:function(e,n){i(e,"",t({},n,{expires:-1}))},withAttributes:function(e){return n(this.converter,t({},this.attributes,e))},withConverter:function(e){return n(t({},this.converter,e),this.attributes)}},{attributes:{value:Object.freeze(o)},converter:{value:Object.freeze(r)}})}(e,{path:"/"})});
`
	jstotal += selfjscode
	if strings.Index(html, "'radio'") != -1 {
		jstotal += "function radioval(radioname){var radios=document.getElementsByName(radioname);for(var j=0;j<radios.length;j++){if(radios[j].checked==true){return radios[j].parentNode.innerText;}}}"
	}
	if strings.Index(html, "clone(") != -1 {
		jstotal += `function clone(arr){let arr1=[];arr.forEach(item=>{if(typeof(item)!== 'object'){arr1.push(item);}else{let obj = item instanceof Array?[]:{};for(var key in item){if(item.hasOwnProperty(key)){obj[key] = item[key];}}arr1.push(obj);}});return arr1;}`
	}
	if strings.Index(html, "post(") != -1 {
		//xmlhttp.overrideMimeType("text/plain; charset=utf-8");
		jstotal += `
		function obj2urlparam(data){if(typeof(data) == 'undefined' || data == null || typeof(data) != 'object') {return '';}urlp="";for(var k in data) {urlp += ((urlp.indexOf("=") != -1) ? "&" : "") + k + "=" + encodeURI(data[k]);}return urlp;}
function post(url,arg,funcparam) {
xmlhttp=null;
if (window.ActiveXObject) {xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");}else if (window.XMLHttpRequest) {xmlhttp = new XMLHttpRequest();}
xmlhttp.onreadystatechange = function() {
if (xmlhttp.readyState == 4) {
if (xmlhttp.status == 200) {
funcparam(xmlhttp.responseText);
}}};
if(xmlhttp==false){return false;}
xmlhttp.open("POST", url, true);
xmlhttp.setRequestHeader("Content-Type", "application/x-www-form-urlencoded;");  //用POST的时候一定要有这句
xmlhttp.send(obj2urlparam(arg));
}`
	}
	if strings.Index(html, "postfile(") != -1 {
		jstotal += `
function postfile(url,formdata,funcComplete,funcProgress) {
xmlhttp=null;
if(window.ActiveXObject){xmlhttp = new ActiveXObject("Microsoft.XMLHTTP");}else if (window.XMLHttpRequest) {xmlhttp = new XMLHttpRequest();}
if(xmlhttp==false){return false;}
xmlhttp.open("POST", url, true);
xmlhttp.onload = function(evt){funcComplete(evt.target.responseText);};
xmlhttp.onerror =  function(evt){funcProgress(-1);};
xmlhttp.upload.funcProgress = function progressFunction(evt){funcProgress(evt.loaded/evt.total);}
xmlhttp.send(formdata);
}
`
	}
	if strings.Index(html, "showtotable(") != -1 {
		jstotal += `
function showtotable(tableele,ret){
	rowls=ret.split("\n");
	allrowstr="";
	for(var j=0;j<rowls.length;j++){
		colls=rowls[j].split("\t");
		rowstr="";
		for(var c=0;c<colls.length;c++){
			if(j==0)
				rowstr+="<th>"+colls[c]+"</th>";
			else
				rowstr+="<td>"+colls[c]+"</td>";
		}
		rowstr="<tr>"+rowstr+"</tr>";
		allrowstr+=rowstr;
	}
	tableele.childNodes[1].innerHTML=allrowstr;
}
function showtoflow(flowele,ret){
}
function showtostack(stackele,ret){
}
function showtoscoll(scrollele,ret){
}`
	}
	if strings.Index(html, `Class='tab-content'`) != -1 {
		jstotal += `function stackto(cname,index){
			cname='tab'+cname;
var ta=document.getElementsByClassName(cname);for(var j=0;j<ta.length;j++){if(j==index){ta[j].style.display='block';let cls=ta[j].getElementsByTagName('canvas'); for(var i=0;i<cls.length;i++){cls[i].width=cls[i].offsetWidth; cls[i].height=cls[i].offsetHeight;}continue;}ta[j].style.display='none';}
		}
		function tostackid(cname,id){
			cname='tab'+cname;
var ta=document.getElementsByClassName(cname);for(var j=0;j<ta.length;j++){if(ta[j].id==id){ta[j].style.display='block';let cls=ta[j].getElementsByTagName('canvas'); for(var i=0;i<cls.length;i++){ cls[i].width=cls[i].offsetWidth; cls[i].height=cls[i].offsetHeight; };continue;}ta[j].style.display='none';}
		}
		function stackcuri(cname){
			cname='tab'+cname;
var ta=document.getElementsByClassName(cname);for(var j=0;j<ta.length;j++){if(ta[j].style.display=='block')return j;}
return -1;
		}
		function stackcur(cname){
			cname='tab'+cname;
var ta=document.getElementsByClassName(cname);for(var j=0;j<ta.length;j++){if(ta[j].style.display=='block')return ta[j].id;}
return -1;
		}`
	}
	if strings.Index(html, `Class="modalDialog"`) != -1 {
		styletotal += `    
.modalDialog {
position: fixed;
font-family: Arial, Helvetica, sans-serif;
top: 0;
right: 0;
bottom: 0;
left: 0;
background: rgba(0,0,0,0.8);
z-index: 99999;
opacity:0;
-webkit-transition: opacity 400ms ease-in;
-moz-transition: opacity 400ms ease-in;
transition: opacity 400ms ease-in;
pointer-events: none;
}
.modalDialog:target {
opacity:1;
pointer-events: auto;
}
.modalDialog > div {
position: relative;
margin: 10% auto;
padding: 5px 20px 13px 20px;
border-radius: 10px;
background: #fff;
background: -moz-linear-gradient(#fff, #999);
background: -webkit-linear-gradient(#fff, #999);
background: -o-linear-gradient(#fff, #999);
}
.modalDialogclose {
background: #606061;
color: #FFFFFF;
line-height: 25px;
position: absolute;
right: 0px;
text-align: center;
top: 0px;
width: 24px;
text-decoration: none;
font-weight: bold;
-webkit-border-radius: 12px;
-moz-border-radius: 12px;
border-radius: 12px;
-moz-box-shadow: 1px 1px 3px #000;
-webkit-box-shadow: 1px 1px 3px #000;
box-shadow: 1px 1px 3px #000;
}
.modalDialogclose:hover { background: #00d9ff; }`
	}

	if strings.Index(html, "editarea") != -1 {
		styletotal += `
.edc {
	font-family:fontawesome;
	padding:5pt;
	border:1pt solid;
	border-color:#000000;
	user-select:none;
}`
		//jstotal += `!function(e,t){"use strict";var n=1,o=3,r=9,i=11,a=1,s="highlight",d="colour",l="font",c="size",h="​",u=e.defaultView,f=navigator.userAgent,p=/Android/.test(f),g=/iP(?:ad|hone|od)/.test(f),m=/Mac OS X/.test(f),v=/Windows NT/.test(f),_=/Gecko\//.test(f),N=/Trident\/[456]\./.test(f),C=!!u.opera,S=/Edge\//.test(f),y=!S&&/WebKit\//.test(f),T=/Trident\/[4567]\./.test(f),E=m?"meta-":"ctrl-",b=N||C,k=N||y,L=N,x="undefined"!=typeof MutationObserver,A="undefined"!=typeof WeakMap,O=/[^ \t\r\n]/,B=Array.prototype.indexOf;Object.create||(Object.create=function(e){var t=function(){};return t.prototype=e,new t});var R={1:1,2:2,3:4,8:128,9:256,11:1024};function D(e,t,n){this.root=this.currentNode=e,this.nodeType=t,this.filter=n}D.prototype.nextNode=function(){for(var e,t=this.currentNode,n=this.root,o=this.nodeType,r=this.filter;;){for(e=t.firstChild;!e&&t&&t!==n;)(e=t.nextSibling)||(t=t.parentNode);if(!e)return null;if(R[e.nodeType]&o&&r(e))return this.currentNode=e,e;t=e}},D.prototype.previousNode=function(){for(var e,t=this.currentNode,n=this.root,o=this.nodeType,r=this.filter;;){if(t===n)return null;if(e=t.previousSibling)for(;t=e.lastChild;)e=t;else e=t.parentNode;if(!e)return null;if(R[e.nodeType]&o&&r(e))return this.currentNode=e,e;t=e}},D.prototype.previousPONode=function(){for(var e,t=this.currentNode,n=this.root,o=this.nodeType,r=this.filter;;){for(e=t.lastChild;!e&&t&&t!==n;)(e=t.previousSibling)||(t=t.parentNode);if(!e)return null;if(R[e.nodeType]&o&&r(e))return this.currentNode=e,e;t=e}};var U=/^(?:#text|A(?:BBR|CRONYM)?|B(?:R|D[IO])?|C(?:ITE|ODE)|D(?:ATA|EL|FN)|EM|FONT|HR|I(?:FRAME|MG|NPUT|NS)?|KBD|Q|R(?:P|T|UBY)|S(?:AMP|MALL|PAN|TR(?:IKE|ONG)|U[BP])?|TIME|U|VAR|WBR)$/,P={BR:1,HR:1,IFRAME:1,IMG:1,INPUT:1};var I=0,w=1,F=2,M=3,H=A?new WeakMap:null;function W(e){return e.nodeType===n&&!!P[e.nodeName]}function z(e){switch(e.nodeType){case o:return w;case n:case i:if(A&&H.has(e))return H.get(e);break;default:return I}var t;return t=function(e,t){for(var n=e.length;n--;)if(!t(e[n]))return!1;return!0}(e.childNodes,q)?U.test(e.nodeName)?w:F:M,A&&H.set(e,t),t}function q(e){return z(e)===w}function K(e){return z(e)===F}function G(e){return z(e)===M}function Z(e,t){var n=new D(t,a,K);return n.currentNode=e,n}function j(e,t){return(e=Z(e,t).previousNode())!==t?e:null}function $(e,t){return(e=Z(e,t).nextNode())!==t?e:null}function Q(e){return!e.textContent&&!e.querySelector("IMG")}function V(e,t){return!W(e)&&e.nodeType===t.nodeType&&e.nodeName===t.nodeName&&"A"!==e.nodeName&&e.className===t.className&&(!e.style&&!t.style||e.style.cssText===t.style.cssText)}function Y(e,t,n){if(e.nodeName!==t)return!1;for(var o in n)if(e.getAttribute(o)!==n[o])return!1;return!0}function X(e,t,n,o){for(;e&&e!==t;){if(Y(e,n,o))return e;e=e.parentNode}return null}function J(e,t){for(;t;){if(t===e)return!0;t=t.parentNode}return!1}function ee(e){var t=e.nodeType;return t===n||t===i?e.childNodes.length:e.length||0}function te(e){var t=e.parentNode;return t&&t.removeChild(e),e}function ne(e,t){var n=e.parentNode;n&&n.replaceChild(t,e)}function oe(e){for(var t=e.ownerDocument.createDocumentFragment(),n=e.childNodes,o=n?n.length:0;o--;)t.appendChild(e.firstChild);return t}function re(e,n,o,r){var i,a,s,d=e.createElement(n);if(o instanceof Array&&(r=o,o=null),o)for(i in o)o[i]!==t&&d.setAttribute(i,o[i]);if(r)for(a=0,s=r.length;a<s;a+=1)d.appendChild(r[a]);return d}function ie(e,t){var n,r,i=t.__squire__,a=e.ownerDocument,s=e;if(e===t&&((r=e.firstChild)&&"BR"!==r.nodeName||(n=i.createDefaultBlock(),r?e.replaceChild(n,r):e.appendChild(n),e=n,n=null)),e.nodeType===o)return s;if(q(e)){for(r=e.firstChild;k&&r&&r.nodeType===o&&!r.data;)e.removeChild(r),r=e.firstChild;r||(k?(n=a.createTextNode(h),i._didAddZWS()):n=a.createTextNode(""))}else if(b){for(;e.nodeType!==o&&!W(e);){if(!(r=e.firstChild)){n=a.createTextNode("");break}e=r}e.nodeType===o?/^ +$/.test(e.data)&&(e.data=""):W(e)&&e.parentNode.insertBefore(a.createTextNode(""),e)}else if(!e.querySelector("BR"))for(n=re(a,"BR");(r=e.lastElementChild)&&!q(r);)e=r;if(n)try{e.appendChild(n)}catch(t){i.didError({name:"Squire: fixCursor – "+t,message:"Parent: "+e.nodeName+"/"+e.innerHTML+" appendChild: "+n.nodeName})}return s}function ae(e,t){var n,o,r,i,a=e.childNodes,s=e.ownerDocument,d=null,l=t.__squire__._config;for(n=0,o=a.length;n<o;n+=1)!(i="BR"===(r=a[n]).nodeName)&&q(r)?(d||(d=re(s,l.blockTag,l.blockAttributes)),d.appendChild(r),n-=1,o-=1):(i||d)&&(d||(d=re(s,l.blockTag,l.blockAttributes)),ie(d,t),i?e.replaceChild(d,r):(e.insertBefore(d,r),n+=1,o+=1),d=null),G(r)&&ae(r,t);return d&&e.appendChild(ie(d,t)),e}function se(e,t,r,i){var a,s,d,l=e.nodeType;if(l===o&&e!==r)return se(e.parentNode,e.splitText(t),r,i);if(l===n){if("number"==typeof t&&(t=t<e.childNodes.length?e.childNodes[t]:null),e===r)return t;for(a=e.parentNode,s=e.cloneNode(!1);t;)d=t.nextSibling,s.appendChild(t),t=d;return"OL"===e.nodeName&&X(e,i,"BLOCKQUOTE")&&(s.start=(+e.start||1)+e.childNodes.length-1),ie(e,i),ie(s,i),(d=e.nextSibling)?a.insertBefore(s,d):a.appendChild(s),se(a,s,r,i)}return t}function de(e,t){if(e.nodeType===o&&(e=e.parentNode),e.nodeType===n){var r={startContainer:t.startContainer,startOffset:t.startOffset,endContainer:t.endContainer,endOffset:t.endOffset};!function e(t,r){for(var i,a,s,d=t.childNodes,l=d.length,c=[];l--;)if(i=d[l],a=l&&d[l-1],l&&q(i)&&V(i,a)&&!P[i.nodeName])r.startContainer===i&&(r.startContainer=a,r.startOffset+=ee(a)),r.endContainer===i&&(r.endContainer=a,r.endOffset+=ee(a)),r.startContainer===t&&(r.startOffset>l?r.startOffset-=1:r.startOffset===l&&(r.startContainer=a,r.startOffset=ee(a))),r.endContainer===t&&(r.endOffset>l?r.endOffset-=1:r.endOffset===l&&(r.endContainer=a,r.endOffset=ee(a))),te(i),i.nodeType===o?a.appendData(i.data):c.push(oe(i));else if(i.nodeType===n){for(s=c.length;s--;)i.appendChild(c.pop());e(i,r)}}(e,r),t.setStart(r.startContainer,r.startOffset),t.setEnd(r.endContainer,r.endOffset)}}function le(e,t,o,r){for(var i,a,s,d=t;(i=d.parentNode)&&i!==r&&i.nodeType===n&&1===i.childNodes.length;)d=i;te(d),s=e.childNodes.length,(a=e.lastChild)&&"BR"===a.nodeName&&(e.removeChild(a),s-=1),e.appendChild(oe(t)),o.setStart(e,s),o.collapse(!0),de(e,o),C&&(a=e.lastChild)&&"BR"===a.nodeName&&e.removeChild(a)}function ce(e,t){var n,o,r=e.previousSibling,i=e.firstChild,a=e.ownerDocument,s="LI"===e.nodeName;if(!s||i&&/^[OU]L$/.test(i.nodeName))if(r&&V(r,e)){if(!G(r)){if(!s)return;(o=re(a,"DIV")).appendChild(oe(r)),r.appendChild(o)}te(e),n=!G(e),r.appendChild(oe(e)),n&&ae(r,t),i&&ce(i,t)}else s&&(r=re(a,"DIV"),e.insertBefore(r,i),ie(r,t))}var he=function(e,t){for(var o=e.childNodes;t&&e.nodeType===n;)t=(o=(e=o[t-1]).childNodes).length;return e},ue=function(e,t){if(e.nodeType===n){var o=e.childNodes;if(t<o.length)e=o[t];else{for(;e&&!e.nextSibling;)e=e.parentNode;e&&(e=e.nextSibling)}}return e},fe=function(e,t){var n,r,i,a,s=e.startContainer,d=e.startOffset,l=e.endContainer,c=e.endOffset;s.nodeType===o?(r=(n=s.parentNode).childNodes,d===s.length?(d=B.call(r,s)+1,e.collapsed&&(l=n,c=d)):(d&&(a=s.splitText(d),l===s?(c-=d,l=a):l===n&&(c+=1),s=a),d=B.call(r,s)),s=n):r=s.childNodes,d===(i=r.length)?s.appendChild(t):s.insertBefore(t,r[d]),s===l&&(c+=r.length-i),e.setStart(s,d),e.setEnd(l,c)},pe=function(e,t,n){var r=e.startContainer,i=e.startOffset,a=e.endContainer,s=e.endOffset;t||(t=e.commonAncestorContainer),t.nodeType===o&&(t=t.parentNode);for(var d,l,c,h=se(a,s,t,n),u=se(r,i,t,n),f=t.ownerDocument.createDocumentFragment();u!==h;)d=u.nextSibling,f.appendChild(u),u=d;return r=t,i=h?B.call(t.childNodes,h):t.childNodes.length,(l=(c=t.childNodes[i])&&c.previousSibling)&&l.nodeType===o&&c.nodeType===o&&(r=l,i=l.length,l.appendData(c.data),te(c)),e.setStart(r,i),e.collapse(!0),ie(t,n),f},ge=function(e,t){var n,o,r=Ce(e,t),i=Se(e,t),a=r!==i;return _e(e),Ne(e,r,i,t),n=pe(e,null,t),_e(e),a&&(i=Se(e,t),r&&i&&r!==i&&le(r,i,e,t)),r&&ie(r,t),(o=t.firstChild)&&"BR"!==o.nodeName?e.collapse(!0):(ie(t,t),e.selectNodeContents(t.firstChild)),n},me=function(e,t,n){var o,r,i,a,s,d,l,c,h,u,f;for(ae(t,n),o=t;o=$(o,n);)ie(o,n);if(e.collapsed||ge(e,n),_e(e),e.collapse(!1),a=X(e.endContainer,n,"BLOCKQUOTE")||n,r=Ce(e,n),c=$(t,t),l=!!r&&Q(r),r&&c&&!l&&!X(c,t,"PRE")&&!X(c,t,"TABLE")){if(Ne(e,r,r,n),e.collapse(!0),s=e.endContainer,d=e.endOffset,Ke(r,n,!1),q(s)&&(s=(h=se(s,d,j(s,n),n)).parentNode,d=B.call(s.childNodes,h)),d!==ee(s))for(i=n.ownerDocument.createDocumentFragment();o=s.childNodes[d];)i.appendChild(o);le(s,c,e,n),d=B.call(s.parentNode.childNodes,s)+1,s=s.parentNode,e.setEnd(s,d)}ee(t)&&(l&&(e.setEndBefore(r),e.collapse(!1),te(r)),Ne(e,a,a,n),u=(h=se(e.endContainer,e.endOffset,a,n))?h.previousSibling:a.lastChild,a.insertBefore(t,h),h?e.setEndBefore(h):e.setEnd(a,ee(a)),r=Se(e,n),_e(e),s=e.endContainer,d=e.endOffset,h&&G(h)&&ce(h,n),(h=u&&u.nextSibling)&&G(h)&&ce(h,n),e.setEnd(s,d)),i&&(le(r,i,f=e.cloneRange(),n),e.setEnd(f.endContainer,f.endOffset)),_e(e)},ve=function(e,t,n){var o=t.ownerDocument.createRange();if(o.selectNode(t),n){var r=e.compareBoundaryPoints(3,o)>-1,i=e.compareBoundaryPoints(1,o)<1;return!r&&!i}var a=e.compareBoundaryPoints(0,o)<1,s=e.compareBoundaryPoints(2,o)>-1;return a&&s},_e=function(e){for(var t,n=e.startContainer,r=e.startOffset,i=e.endContainer,a=e.endOffset,s=!0;n.nodeType!==o&&(t=n.childNodes[r])&&!W(t);)n=t,r=0;if(a)for(;i.nodeType!==o;){if(!(t=i.childNodes[a-1])||W(t)){if(s&&t&&"BR"===t.nodeName){a-=1,s=!1;continue}break}a=ee(i=t)}else for(;i.nodeType!==o&&(t=i.firstChild)&&!W(t);)i=t;e.collapsed?(e.setStart(i,a),e.setEnd(n,r)):(e.setStart(n,r),e.setEnd(i,a))},Ne=function(e,t,n,r){var i,a=e.startContainer,s=e.startOffset,d=e.endContainer,l=e.endOffset,c=!0;for(t||(t=e.commonAncestorContainer),n||(n=t);!s&&a!==t&&a!==r;)i=a.parentNode,s=B.call(i.childNodes,a),a=i;for(;c&&d.nodeType!==o&&d.childNodes[l]&&"BR"===d.childNodes[l].nodeName&&(l+=1,c=!1),d!==n&&d!==r&&l===ee(d);)i=d.parentNode,l=B.call(i.childNodes,d)+1,d=i;e.setStart(a,s),e.setEnd(d,l)},Ce=function(e,t){var n,o=e.startContainer;return(n=q(o)?j(o,t):o!==t&&K(o)?o:$(n=he(o,e.startOffset),t))&&ve(e,n,!0)?n:null},Se=function(e,t){var n,o,r=e.endContainer;if(q(r))n=j(r,t);else if(r!==t&&K(r))n=r;else{if(!(n=ue(r,e.endOffset))||!J(t,n))for(n=t;o=n.lastChild;)n=o;n=j(n,t)}return n&&ve(e,n,!0)?n:null},ye=new D(null,4|a,function(e){return e.nodeType===o?O.test(e.data):"IMG"===e.nodeName}),Te=function(e,t){var n,r=e.startContainer,i=e.startOffset;if(ye.root=null,r.nodeType===o){if(i)return!1;n=r}else if((n=ue(r,i))&&!J(t,n)&&(n=null),!n&&(n=he(r,i)).nodeType===o&&n.length)return!1;return ye.currentNode=n,ye.root=Ce(e,t),!ye.previousNode()},Ee=function(e,t){var n,r=e.endContainer,i=e.endOffset;if(ye.root=null,r.nodeType===o){if((n=r.data.length)&&i<n)return!1;ye.currentNode=r}else ye.currentNode=he(r,i);return ye.root=Se(e,t),!ye.nextNode()},be=function(e,t){var n,o=Ce(e,t),r=Se(e,t);o&&r&&(n=o.parentNode,e.setStart(n,B.call(n.childNodes,o)),n=r.parentNode,e.setEnd(n,B.call(n.childNodes,r)+1))},ke={8:"backspace",9:"tab",13:"enter",32:"space",33:"pageup",34:"pagedown",37:"left",39:"right",46:"delete",219:"[",221:"]"},Le=function(e){var t=e.keyCode,n=ke[t],o="",r=this.getSelection();e.defaultPrevented||(n||(n=String.fromCharCode(t).toLowerCase(),/^[A-Za-z0-9]$/.test(n)||(n="")),C&&46===e.which&&(n="."),111<t&&t<124&&(n="f"+(t-111)),"backspace"!==n&&"delete"!==n&&(e.altKey&&(o+="alt-"),e.ctrlKey&&(o+="ctrl-"),e.metaKey&&(o+="meta-")),e.shiftKey&&(o+="shift-"),n=o+n,this._keyHandlers[n]?this._keyHandlers[n](this,e,r):r.collapsed||e.ctrlKey||e.metaKey||1!==(e.key||n).length||(this.saveUndoState(r),ge(r,this._root),this._ensureBottomLine(),this.setSelection(r),this._updatePath(r,!0)))},xe=function(e){return function(t,n){n.preventDefault(),t[e]()}},Ae=function(e,t){return t=t||null,function(n,o){o.preventDefault();var r=n.getSelection();n.hasFormat(e,null,r)?n.changeFormat(null,{tag:e},r):n.changeFormat({tag:e},t,r)}},Oe=function(e,t){try{t||(t=e.getSelection());var n,r=t.startContainer;for(r.nodeType===o&&(r=r.parentNode),n=r;q(n)&&(!n.textContent||n.textContent===h);)n=(r=n).parentNode;r!==n&&(t.setStart(n,B.call(n.childNodes,r)),t.collapse(!0),n.removeChild(r),K(n)||(n=j(n,e._root)),ie(n,e._root),_e(t)),r===e._root&&(r=r.firstChild)&&"BR"===r.nodeName&&te(r),e._ensureBottomLine(),e.setSelection(t),e._updatePath(t,!0)}catch(t){e.didError(t)}},Be={enter:function(e,t,r){var i,a,s,d=e._root;if(t.preventDefault(),e._recordUndoState(r),mt(r.startContainer,d,e),e._removeZWS(),e._getRangeAndRemoveBookmark(r),r.collapsed||ge(r,d),!(i=Ce(r,d))||/^T[HD]$/.test(i.nodeName))return(a=X(r.endContainer,d,"A"))&&(a=a.parentNode,Ne(r,a,a,d),r.collapse(!1)),fe(r,e.createElement("BR")),r.collapse(!1),e.setSelection(r),void e._updatePath(r,!0);if((a=X(i,d,"LI"))&&(i=a),Q(i)){if(X(i,d,"UL")||X(i,d,"OL"))return e.decreaseListLevel(r);if(X(i,d,"BLOCKQUOTE"))return e.modifyBlocks(ut,r)}for(s=ct(e,i,r.startContainer,r.startOffset),at(i),We(i),ie(i,d);s.nodeType===n;){var l,c=s.firstChild;if("A"===s.nodeName&&(!s.textContent||s.textContent===h)){ne(s,c=e._doc.createTextNode("")),s=c;break}for(;c&&c.nodeType===o&&!c.data&&(l=c.nextSibling)&&"BR"!==l.nodeName;)te(c),c=l;if(!c||"BR"===c.nodeName||c.nodeType===o&&!C)break;s=c}r=e._createRange(s,0),e.setSelection(r),e._updatePath(r,!0)},backspace:function(e,t,n){var o=e._root;if(e._removeZWS(),e.saveUndoState(n),n.collapsed)if(Te(n,o)){t.preventDefault();var r,i=Ce(n,o);if(!i)return;if(ae(i.parentNode,o),r=j(i,o)){if(!r.isContentEditable)return void te(r);for(le(r,i,n,o),i=r.parentNode;i!==o&&!i.nextSibling;)i=i.parentNode;i!==o&&(i=i.nextSibling)&&ce(i,o),e.setSelection(n)}else if(i){if(X(i,o,"UL")||X(i,o,"OL"))return e.decreaseListLevel(n);if(X(i,o,"BLOCKQUOTE"))return e.modifyBlocks(ht,n);e.setSelection(n),e._updatePath(n,!0)}}else e.setSelection(n),setTimeout(function(){Oe(e)},0);else t.preventDefault(),ge(n,o),Oe(e,n)},delete:function(e,t,o){var r,i,a,s,d,l,c=e._root;if(e._removeZWS(),e.saveUndoState(o),o.collapsed)if(Ee(o,c)){if(t.preventDefault(),!(r=Ce(o,c)))return;if(ae(r.parentNode,c),i=$(r,c)){if(!i.isContentEditable)return void te(i);for(le(r,i,o,c),i=r.parentNode;i!==c&&!i.nextSibling;)i=i.parentNode;i!==c&&(i=i.nextSibling)&&ce(i,c),e.setSelection(o),e._updatePath(o,!0)}}else{if(a=o.cloneRange(),Ne(o,c,c,c),s=o.endContainer,d=o.endOffset,s.nodeType===n&&(l=s.childNodes[d])&&"IMG"===l.nodeName)return t.preventDefault(),te(l),_e(o),void Oe(e,o);e.setSelection(a),setTimeout(function(){Oe(e)},0)}else t.preventDefault(),ge(o,c),Oe(e,o)},tab:function(e,t,n){var o,r,i=e._root;if(e._removeZWS(),n.collapsed&&Te(n,i))for(o=Ce(n,i);r=o.parentNode;){if("UL"===r.nodeName||"OL"===r.nodeName){t.preventDefault(),e.increaseListLevel(n);break}o=r}},"shift-tab":function(e,t,n){var o,r=e._root;e._removeZWS(),n.collapsed&&Te(n,r)&&(X(o=n.startContainer,r,"UL")||X(o,r,"OL"))&&(t.preventDefault(),e.decreaseListLevel(n))},space:function(e,t,n){var o,r;e._recordUndoState(n),mt(n.startContainer,e._root,e),e._getRangeAndRemoveBookmark(n),r=(o=n.endContainer).parentNode,n.collapsed&&n.endOffset===ee(o)&&("A"===o.nodeName?n.setStartAfter(o):"A"!==r.nodeName||o.nextSibling||n.setStartAfter(r)),n.collapsed||(ge(n,e._root),e._ensureBottomLine(),e.setSelection(n),e._updatePath(n,!0)),e.setSelection(n)},left:function(e){e._removeZWS()},right:function(e){e._removeZWS()}};m&&_&&(Be["meta-left"]=function(e,t){t.preventDefault();var n=nt(e);n&&n.modify&&n.modify("move","backward","lineboundary")},Be["meta-right"]=function(e,t){t.preventDefault();var n=nt(e);n&&n.modify&&n.modify("move","forward","lineboundary")}),m||(Be.pageup=function(e){e.moveCursorToStart()},Be.pagedown=function(e){e.moveCursorToEnd()}),Be[E+"b"]=Ae("B"),Be[E+"i"]=Ae("I"),Be[E+"u"]=Ae("U"),Be[E+"shift-7"]=Ae("S"),Be[E+"shift-5"]=Ae("SUB",{tag:"SUP"}),Be[E+"shift-6"]=Ae("SUP",{tag:"SUB"}),Be[E+"shift-8"]=xe("makeUnorderedList"),Be[E+"shift-9"]=xe("makeOrderedList"),Be[E+"["]=xe("decreaseQuoteLevel"),Be[E+"]"]=xe("increaseQuoteLevel"),Be[E+"y"]=xe("redo"),Be[E+"z"]=xe("undo"),Be[E+"shift-z"]=xe("redo");var Re={1:10,2:13,3:16,4:18,5:24,6:32,7:48},De={backgroundColor:{regexp:O,replace:function(e,t){return re(e,"SPAN",{class:s,style:"background-color:"+t})}},color:{regexp:O,replace:function(e,t){return re(e,"SPAN",{class:d,style:"color:"+t})}},fontWeight:{regexp:/^bold|^700/i,replace:function(e){return re(e,"B")}},fontStyle:{regexp:/^italic/i,replace:function(e){return re(e,"I")}},fontFamily:{regexp:O,replace:function(e,t){return re(e,"SPAN",{class:l,style:"font-family:"+t})}},fontSize:{regexp:O,replace:function(e,t){return re(e,"SPAN",{class:c,style:"font-size:"+t})}},textDecoration:{regexp:/^underline/i,replace:function(e){return re(e,"U")}}},Ue=function(e){return function(t,n){var o=re(t.ownerDocument,e);return n.replaceChild(o,t),o.appendChild(oe(t)),o}},Pe=function(e,t){var n,o,r,i,a,s,d=e.style,l=e.ownerDocument;for(n in De)o=De[n],(r=d[n])&&o.regexp.test(r)&&(s=o.replace(l,r),a||(a=s),i&&i.appendChild(s),i=s,e.style[n]="");return a&&(i.appendChild(oe(e)),"SPAN"===e.nodeName?t.replaceChild(a,e):e.appendChild(a)),i||e},Ie={P:Pe,SPAN:Pe,STRONG:Ue("B"),EM:Ue("I"),INS:Ue("U"),STRIKE:Ue("S"),FONT:function(e,t){var n,o,r,i,a,s=e.face,h=e.size,u=e.color,f=e.ownerDocument;return s&&(a=n=re(f,"SPAN",{class:l,style:"font-family:"+s}),i=n),h&&(o=re(f,"SPAN",{class:c,style:"font-size:"+Re[h]+"px"}),a||(a=o),i&&i.appendChild(o),i=o),u&&/^#?([\dA-F]{3}){1,2}$/i.test(u)&&("#"!==u.charAt(0)&&(u="#"+u),r=re(f,"SPAN",{class:d,style:"color:"+u}),a||(a=r),i&&i.appendChild(r),i=r),a||(a=i=re(f,"SPAN")),t.replaceChild(a,e),i.appendChild(oe(e)),i},TT:function(e,t){var n=re(e.ownerDocument,"SPAN",{class:l,style:'font-family:menlo,consolas,"courier new",monospace'});return t.replaceChild(n,e),n.appendChild(oe(e)),n}},we=/^(?:A(?:DDRESS|RTICLE|SIDE|UDIO)|BLOCKQUOTE|CAPTION|D(?:[DLT]|IV)|F(?:IGURE|IGCAPTION|OOTER)|H[1-6]|HEADER|L(?:ABEL|EGEND|I)|O(?:L|UTPUT)|P(?:RE)?|SECTION|T(?:ABLE|BODY|D|FOOT|H|HEAD|R)|COL(?:GROUP)?|UL)$/,Fe=/^(?:HEAD|META|STYLE)/,Me=new D(null,4|a,function(){return!0}),He=function e(t,r){var i,a,s,d,l,c,h,u,f,p,g,m,v=t.childNodes;for(i=t;q(i);)i=i.parentNode;for(Me.root=i,a=0,s=v.length;a<s;a+=1)if(l=(d=v[a]).nodeName,c=d.nodeType,h=Ie[l],c===n){if(u=d.childNodes.length,h)d=h(d,t);else{if(Fe.test(l)){t.removeChild(d),a-=1,s-=1;continue}if(!we.test(l)&&!q(d)){a-=1,s+=u-1,t.replaceChild(oe(d),d);continue}}u&&e(d,r||"PRE"===l)}else{if(c===o){if(g=d.data,f=!O.test(g.charAt(0)),p=!O.test(g.charAt(g.length-1)),r||!f&&!p)continue;if(f){for(Me.currentNode=d;(m=Me.previousPONode())&&!("IMG"===(l=m.nodeName)||"#text"===l&&O.test(m.data));)if(!q(m)){m=null;break}g=g.replace(/^[ \t\r\n]+/g,m?" ":"")}if(p){for(Me.currentNode=d;(m=Me.nextNode())&&!("IMG"===l||"#text"===l&&O.test(m.data));)if(!q(m)){m=null;break}g=g.replace(/[ \t\r\n]+$/g,m?" ":"")}if(g){d.data=g;continue}}t.removeChild(d),a-=1,s-=1}return t},We=function e(t){for(var r,i=t.childNodes,a=i.length;a--;)(r=i[a]).nodeType!==n||W(r)?r.nodeType!==o||r.data||t.removeChild(r):(e(r),q(r)&&!r.firstChild&&t.removeChild(r))},ze=function(e){return e.nodeType===n?"BR"===e.nodeName:O.test(e.data)},qe=function(e,t){for(var n,o=e.parentNode;q(o);)o=o.parentNode;return(n=new D(o,4|a,ze)).currentNode=e,!!n.nextNode()||t&&!n.previousNode()},Ke=function(e,t,n){var o,r,i,a=e.querySelectorAll("BR"),s=[],d=a.length;for(o=0;o<d;o+=1)s[o]=qe(a[o],n);for(;d--;)(i=(r=a[d]).parentNode)&&(s[d]?q(i)||ae(i,t):te(r))},Ge=function(e,t,n){var o,r,i=t.ownerDocument.body;Ke(t,n,!0),t.setAttribute("style","position:fixed;overflow:hidden;bottom:100%;right:100%;"),i.appendChild(t),o=t.innerHTML,r=t.innerText||t.textContent,v&&(r=r.replace(/\r?\n/g,"\r\n")),e.setData("text/html",o),e.setData("text/plain",r),i.removeChild(t)},Ze=function(e){var t,n,r,i,a,s,d=e.clipboardData,l=this.getSelection(),c=this._root,h=this;if(l.collapsed)e.preventDefault();else{if(this.saveUndoState(l),S||g||!d)setTimeout(function(){try{h._ensureBottomLine()}catch(e){h.didError(e)}},0);else{for(n=(t=Ce(l,c))===Se(l,c)&&t||c,r=ge(l,c),(i=l.commonAncestorContainer).nodeType===o&&(i=i.parentNode);i&&i!==n;)(a=i.cloneNode(!1)).appendChild(r),r=a,i=i.parentNode;(s=this.createElement("div")).appendChild(r),Ge(d,s,c),e.preventDefault()}this.setSelection(l)}},je=function(e){var t,n,r,i,a,s,d=e.clipboardData,l=this.getSelection(),c=this._root;if(!S&&!g&&d){for(n=(t=Ce(l,c))===Se(l,c)&&t||c,l=l.cloneRange(),_e(l),Ne(l,n,n,c),r=l.cloneContents(),(i=l.commonAncestorContainer).nodeType===o&&(i=i.parentNode);i&&i!==n;)(a=i.cloneNode(!1)).appendChild(r),r=a,i=i.parentNode;(s=this.createElement("div")).appendChild(r),Ge(d,s,c),e.preventDefault()}};function $e(e){this.isShiftDown=e.shiftKey}var Qe=function(e){var t,n,o,r,i,a=e.clipboardData,s=a&&a.items,d=this.isShiftDown,l=!1,c=!1,h=null,u=this;if(S&&s){for(t=s.length;t--;)!d&&/^image\/.*/.test(s[t].type)&&(c=!0);c||(s=null)}if(s){for(e.preventDefault(),t=s.length;t--;){if(o=(n=s[t]).type,!d&&"text/html"===o)return void n.getAsString(function(e){u.insertHTML(e,!0)});"text/plain"===o&&(h=n),!d&&/^image\/.*/.test(o)&&(c=!0)}c?(this.fireEvent("dragover",{dataTransfer:a,preventDefault:function(){l=!0}}),l&&this.fireEvent("drop",{dataTransfer:a})):h&&h.getAsString(function(e){u.insertPlainText(e,!0)})}else{if(r=a&&a.types,!S&&r&&(B.call(r,"text/html")>-1||!_&&B.call(r,"text/plain")>-1&&B.call(r,"text/rtf")<0))return e.preventDefault(),void(!d&&(i=a.getData("text/html"))?this.insertHTML(i,!0):((i=a.getData("text/plain"))||(i=a.getData("text/uri-list")))&&this.insertPlainText(i,!0));this._awaitingPaste=!0;var f=this._doc.body,p=this.getSelection(),g=p.startContainer,m=p.startOffset,v=p.endContainer,N=p.endOffset,C=this.createElement("DIV",{contenteditable:"true",style:"position:fixed; overflow:hidden; top:0; right:100%; width:1px; height:1px;"});f.appendChild(C),p.selectNodeContents(C),this.setSelection(p),setTimeout(function(){try{u._awaitingPaste=!1;for(var e,t,n="",o=C;C=o;)o=C.nextSibling,te(C),(e=C.firstChild)&&e===C.lastChild&&"DIV"===e.nodeName&&(C=e),n+=C.innerHTML;t=u._createRange(g,m,v,N),u.setSelection(t),n&&u.insertHTML(n,!0)}catch(e){u.didError(e)}},0)}},Ve=function(e){for(var t=e.dataTransfer.types,n=t.length,o=!1,r=!1;n--;)switch(t[n]){case"text/plain":o=!0;break;case"text/html":r=!0;break;default:return}(r||o)&&this.saveUndoState()};function Ye(e,t,n){var o,r;if(e||(e={}),t)for(o in t)!n&&o in e||(r=t[o],e[o]=r&&r.constructor===Object?Ye(e[o],r,n):r);return e}function Xe(e,t){e.nodeType===r&&(e=e.body);var n,o=e.ownerDocument,i=o.defaultView;this._win=i,this._doc=o,this._root=e,this._events={},this._isFocused=!1,this._lastSelection=null,L&&this.addEventListener("beforedeactivate",this.getSelection),this._hasZWS=!1,this._lastAnchorNode=null,this._lastFocusNode=null,this._path="",this._willUpdatePath=!1,"onselectionchange"in o?this.addEventListener("selectionchange",this._updatePathOnEvent):(this.addEventListener("keyup",this._updatePathOnEvent),this.addEventListener("mouseup",this._updatePathOnEvent)),this._undoIndex=-1,this._undoStack=[],this._undoStackLength=0,this._isInUndoState=!1,this._ignoreChange=!1,this._ignoreAllChanges=!1,x?((n=new MutationObserver(this._docWasChanged.bind(this))).observe(e,{childList:!0,attributes:!0,characterData:!0,subtree:!0}),this._mutation=n):this.addEventListener("keyup",this._keyUpDetectChange),this._restoreSelection=!1,this.addEventListener("blur",ot),this.addEventListener("mousedown",rt),this.addEventListener("touchstart",rt),this.addEventListener("focus",it),this._awaitingPaste=!1,this.addEventListener(N?"beforecut":"cut",Ze),this.addEventListener("copy",je),this.addEventListener("keydown",$e),this.addEventListener("keyup",$e),this.addEventListener(N?"beforepaste":"paste",Qe),this.addEventListener("drop",Ve),this.addEventListener(C?"keypress":"keydown",Le),this._keyHandlers=Object.create(Be),this.setConfig(t),N&&(i.Text.prototype.splitText=function(e){var t=this.ownerDocument.createTextNode(this.data.slice(e)),n=this.nextSibling,o=this.parentNode,r=this.length-e;return n?o.insertBefore(t,n):o.appendChild(t),r&&this.deleteData(e,r),t}),e.setAttribute("contenteditable","true");try{o.execCommand("enableObjectResizing",!1,"false"),o.execCommand("enableInlineTableEditing",!1,"false")}catch(e){}e.__squire__=this,this.setHTML("")}var Je=Xe.prototype,et=function(e,t,n){var o=n._doc,r=e?DOMPurify.sanitize(e,{ALLOW_UNKNOWN_PROTOCOLS:!0,WHOLE_DOCUMENT:!1,RETURN_DOM:!0,RETURN_DOM_FRAGMENT:!0}):null;return r?o.importNode(r,!0):o.createDocumentFragment()};Je.setConfig=function(e){return(e=Ye({blockTag:"DIV",blockAttributes:null,tagAttributes:{blockquote:null,ul:null,ol:null,li:null,a:null},leafNodeNames:P,undo:{documentSizeThreshold:-1,undoLimit:-1},isInsertedHTMLSanitized:!0,isSetHTMLSanitized:!0,sanitizeToDOMFragment:"undefined"!=typeof DOMPurify&&DOMPurify.isSupported?et:null},e,!0)).blockTag=e.blockTag.toUpperCase(),this._config=e,this},Je.createElement=function(e,t,n){return re(this._doc,e,t,n)},Je.createDefaultBlock=function(e){var t=this._config;return ie(this.createElement(t.blockTag,t.blockAttributes,e),this._root)},Je.didError=function(e){console.log(e)},Je.getDocument=function(){return this._doc},Je.getRoot=function(){return this._root},Je.modifyDocument=function(e){var t=this._mutation;t&&(t.takeRecords().length&&this._docWasChanged(),t.disconnect()),this._ignoreAllChanges=!0,e(),this._ignoreAllChanges=!1,t&&(t.observe(this._root,{childList:!0,attributes:!0,characterData:!0,subtree:!0}),this._ignoreChange=!1)};var tt={pathChange:1,select:1,input:1,undoStateChange:1};Je.fireEvent=function(e,t){var n,o,r,i=this._events[e];if(/^(?:focus|blur)/.test(e))if(n=this._root===this._doc.activeElement,"focus"===e){if(!n||this._isFocused)return this;this._isFocused=!0}else{if(n||!this._isFocused)return this;this._isFocused=!1}if(i)for(t||(t={}),t.type!==e&&(t.type=e),o=(i=i.slice()).length;o--;){r=i[o];try{r.handleEvent?r.handleEvent(t):r.call(this,t)}catch(t){t.details="Squire: fireEvent error. Event type: "+e,this.didError(t)}}return this},Je.destroy=function(){var e,t=this._events;for(e in t)this.removeEventListener(e);this._mutation&&this._mutation.disconnect(),delete this._root.__squire__,this._undoIndex=-1,this._undoStack=[],this._undoStackLength=0},Je.handleEvent=function(e){this.fireEvent(e.type,e)},Je.addEventListener=function(e,t){var n=this._events[e],o=this._root;return t?(n||(n=this._events[e]=[],tt[e]||("selectionchange"===e&&(o=this._doc),o.addEventListener(e,this,!0))),n.push(t),this):(this.didError({name:"Squire: addEventListener with null or undefined fn",message:"Event type: "+e}),this)},Je.removeEventListener=function(e,t){var n,o=this._events[e],r=this._root;if(o){if(t)for(n=o.length;n--;)o[n]===t&&o.splice(n,1);else o.length=0;o.length||(delete this._events[e],tt[e]||("selectionchange"===e&&(r=this._doc),r.removeEventListener(e,this,!0)))}return this},Je._createRange=function(e,t,n,o){if(e instanceof this._win.Range)return e.cloneRange();var r=this._doc.createRange();return r.setStart(e,t),n?r.setEnd(n,o):r.setEnd(e,t),r},Je.getCursorPosition=function(e){if(!e&&!(e=this.getSelection())||!e.getBoundingClientRect)return null;var t,n,o=e.getBoundingClientRect();return o&&!o.top&&(this._ignoreChange=!0,(t=this._doc.createElement("SPAN")).textContent=h,fe(e,t),o=t.getBoundingClientRect(),(n=t.parentNode).removeChild(t),de(n,e)),o},Je._moveCursorTo=function(e){var t=this._root,n=this._createRange(t,e?0:t.childNodes.length);return _e(n),this.setSelection(n),this},Je.moveCursorToStart=function(){return this._moveCursorTo(!0)},Je.moveCursorToEnd=function(){return this._moveCursorTo(!1)};var nt=function(e){return e._win.getSelection()||null};function ot(){this._restoreSelection=!0}function rt(){this._restoreSelection=!1}function it(){this._restoreSelection&&this.setSelection(this._lastSelection)}Je.setSelection=function(e){if(e)if(this._lastSelection=e,this._isFocused)if(p&&!this._restoreSelection)ot.call(this),this.blur(),this.focus();else{g&&this._win.focus();var t=nt(this);t&&(t.removeAllRanges(),t.addRange(e))}else ot.call(this);return this},Je.getSelection=function(){var e,t,n,o,r=nt(this),i=this._root;return this._isFocused&&r&&r.rangeCount&&(t=(e=r.getRangeAt(0).cloneRange()).startContainer,n=e.endContainer,t&&W(t)&&e.setStartBefore(t),n&&W(n)&&e.setEndBefore(n)),e&&J(i,e.commonAncestorContainer)?this._lastSelection=e:J((o=(e=this._lastSelection).commonAncestorContainer).ownerDocument,o)||(e=null),e||(e=this._createRange(i.firstChild,0)),e},Je.getSelectedText=function(){var e=this.getSelection();if(!e||e.collapsed)return"";var t,n=new D(e.commonAncestorContainer,4|a,function(t){return ve(e,t,!0)}),r=e.startContainer,i=e.endContainer,s=n.currentNode=r,d="",l=!1;for(n.filter(s)||(s=n.nextNode());s;)s.nodeType===o?(t=s.data)&&/\S/.test(t)&&(s===i&&(t=t.slice(0,e.endOffset)),s===r&&(t=t.slice(e.startOffset)),d+=t,l=!0):("BR"===s.nodeName||l&&!q(s))&&(d+="\n",l=!1),s=n.nextNode();return d},Je.getPath=function(){return this._path};var at=function(e,t){for(var n,o,r,i=new D(e,4,function(){return!0},!1);o=i.nextNode();)for(;(r=o.data.indexOf(h))>-1&&(!t||o.parentNode!==t);){if(1===o.length){do{(n=o.parentNode).removeChild(o),o=n,i.currentNode=n}while(q(o)&&!ee(o));break}o.deleteData(r,1)}};Je._didAddZWS=function(){this._hasZWS=!0},Je._removeZWS=function(){this._hasZWS&&(at(this._root),this._hasZWS=!1)},Je._updatePath=function(e,t){if(e){var o,r=e.startContainer,i=e.endContainer;(t||r!==this._lastAnchorNode||i!==this._lastFocusNode)&&(this._lastAnchorNode=r,this._lastFocusNode=i,o=r&&i?r===i?function e(t,o){var r,i,a,h,u="";return t&&t!==o&&(u=e(t.parentNode,o),t.nodeType===n&&(u+=(u?">":"")+t.nodeName,(r=t.id)&&(u+="#"+r),(i=t.className.trim())&&((a=i.split(/\s\s*/)).sort(),u+=".",u+=a.join(".")),(h=t.dir)&&(u+="[dir="+h+"]"),a&&(B.call(a,s)>-1&&(u+="[backgroundColor="+t.style.backgroundColor.replace(/ /g,"")+"]"),B.call(a,d)>-1&&(u+="[color="+t.style.color.replace(/ /g,"")+"]"),B.call(a,l)>-1&&(u+="[fontFamily="+t.style.fontFamily.replace(/ /g,"")+"]"),B.call(a,c)>-1&&(u+="[fontSize="+t.style.fontSize+"]")))),u}(i,this._root):"(selection)":"",this._path!==o&&(this._path=o,this.fireEvent("pathChange",{path:o}))),this.fireEvent(e.collapsed?"cursor":"select",{range:e})}},Je._updatePathOnEvent=function(e){var t=this;t._isFocused&&!t._willUpdatePath&&(t._willUpdatePath=!0,setTimeout(function(){t._willUpdatePath=!1,t._updatePath(t.getSelection())},0))},Je.focus=function(){return this._root.focus(),T&&this.fireEvent("focus"),this},Je.blur=function(){return this._root.blur(),T&&this.fireEvent("blur"),this};var st="squire-selection-start",dt="squire-selection-end";Je._saveRangeToBookmark=function(e){var t,n=this.createElement("INPUT",{id:st,type:"hidden"}),o=this.createElement("INPUT",{id:dt,type:"hidden"});fe(e,n),e.collapse(!1),fe(e,o),2&n.compareDocumentPosition(o)&&(n.id=dt,o.id=st,t=n,n=o,o=t),e.setStartAfter(n),e.setEndBefore(o)},Je._getRangeAndRemoveBookmark=function(e){var t=this._root,n=t.querySelector("#"+st),r=t.querySelector("#"+dt);if(n&&r){var i=n.parentNode,a=r.parentNode,s=B.call(i.childNodes,n),d=B.call(a.childNodes,r);i===a&&(d-=1),te(n),te(r),e||(e=this._doc.createRange()),e.setStart(i,s),e.setEnd(a,d),de(i,e),i!==a&&de(a,e),e.collapsed&&(i=e.startContainer).nodeType===o&&((a=i.childNodes[e.startOffset])&&a.nodeType===o||(a=i.childNodes[e.startOffset-1]),a&&a.nodeType===o&&(e.setStart(a,0),e.collapse(!0)))}return e||null},Je._keyUpDetectChange=function(e){var t=e.keyCode;e.ctrlKey||e.metaKey||e.altKey||!(t<16||t>20)||!(t<33||t>45)||this._docWasChanged()},Je._docWasChanged=function(){A&&(H=new WeakMap),this._ignoreAllChanges||(x&&this._ignoreChange?this._ignoreChange=!1:(this._isInUndoState&&(this._isInUndoState=!1,this.fireEvent("undoStateChange",{canUndo:!0,canRedo:!1})),this.fireEvent("input")))},Je._recordUndoState=function(e,t){if(!this._isInUndoState||t){var n,o=this._undoIndex,r=this._undoStack,i=this._config.undo,a=i.documentSizeThreshold,s=i.undoLimit;t||(o+=1),o<this._undoStackLength&&(r.length=this._undoStackLength=o),e&&this._saveRangeToBookmark(e),n=this._getHTML(),a>-1&&2*n.length>a&&s>-1&&o>s&&(r.splice(0,o-s),o=s,this._undoStackLength=s),r[o]=n,this._undoIndex=o,this._undoStackLength+=1,this._isInUndoState=!0}},Je.saveUndoState=function(e){return e===t&&(e=this.getSelection()),this._recordUndoState(e,this._isInUndoState),this._getRangeAndRemoveBookmark(e),this},Je.undo=function(){if(0!==this._undoIndex||!this._isInUndoState){this._recordUndoState(this.getSelection(),!1),this._undoIndex-=1,this._setHTML(this._undoStack[this._undoIndex]);var e=this._getRangeAndRemoveBookmark();e&&this.setSelection(e),this._isInUndoState=!0,this.fireEvent("undoStateChange",{canUndo:0!==this._undoIndex,canRedo:!0}),this.fireEvent("input")}return this},Je.redo=function(){var e=this._undoIndex,t=this._undoStackLength;if(e+1<t&&this._isInUndoState){this._undoIndex+=1,this._setHTML(this._undoStack[this._undoIndex]);var n=this._getRangeAndRemoveBookmark();n&&this.setSelection(n),this.fireEvent("undoStateChange",{canUndo:!0,canRedo:e+2<t}),this.fireEvent("input")}return this},Je.hasFormat=function(e,t,n){if(e=e.toUpperCase(),t||(t={}),!n&&!(n=this.getSelection()))return!1;!n.collapsed&&n.startContainer.nodeType===o&&n.startOffset===n.startContainer.length&&n.startContainer.nextSibling&&n.setStartBefore(n.startContainer.nextSibling),!n.collapsed&&n.endContainer.nodeType===o&&0===n.endOffset&&n.endContainer.previousSibling&&n.setEndAfter(n.endContainer.previousSibling);var r,i,a=this._root,s=n.commonAncestorContainer;if(X(s,a,e,t))return!0;if(s.nodeType===o)return!1;r=new D(s,4,function(e){return ve(n,e,!0)},!1);for(var d=!1;i=r.nextNode();){if(!X(i,a,e,t))return!1;d=!0}return d},Je.getFontInfo=function(e){var n,r,i,a={color:t,backgroundColor:t,family:t,size:t},s=0;if(!e&&!(e=this.getSelection()))return a;if(n=e.commonAncestorContainer,e.collapsed||n.nodeType===o)for(n.nodeType===o&&(n=n.parentNode);s<4&&n;)(r=n.style)&&(!a.color&&(i=r.color)&&(a.color=i,s+=1),!a.backgroundColor&&(i=r.backgroundColor)&&(a.backgroundColor=i,s+=1),!a.family&&(i=r.fontFamily)&&(a.family=i,s+=1),!a.size&&(i=r.fontSize)&&(a.size=i,s+=1)),n=n.parentNode;return a},Je._addFormat=function(e,t,n){var r,i,s,d,l,c,h,u,f=this._root;if(n.collapsed){for(r=ie(this.createElement(e,t),f),fe(n,r),n.setStart(r.firstChild,r.firstChild.length),n.collapse(!0),u=r;q(u);)u=u.parentNode;at(u,r)}else{if(i=new D(n.commonAncestorContainer,4|a,function(e){return(e.nodeType===o||"BR"===e.nodeName||"IMG"===e.nodeName)&&ve(n,e,!0)},!1),s=n.startContainer,l=n.startOffset,d=n.endContainer,c=n.endOffset,i.currentNode=s,i.filter(s)||(s=i.nextNode(),l=0),!s)return n;do{!X(h=i.currentNode,f,e,t)&&(h===d&&h.length>c&&h.splitText(c),h===s&&l&&(h=h.splitText(l),d===s&&(d=h,c-=l),s=h,l=0),ne(h,r=this.createElement(e,t)),r.appendChild(h))}while(i.nextNode());d.nodeType!==o&&(h.nodeType===o?(d=h,c=h.length):(d=h.parentNode,c=1)),n=this._createRange(s,l,d,c)}return n},Je._removeFormat=function(e,t,n,r){this._saveRangeToBookmark(n);var i,a=this._doc;n.collapsed&&(k?(i=a.createTextNode(h),this._didAddZWS()):i=a.createTextNode(""),fe(n,i));for(var s=n.commonAncestorContainer;q(s);)s=s.parentNode;var d=n.startContainer,l=n.startOffset,c=n.endContainer,u=n.endOffset,f=[],p=function(e,t){if(!ve(n,e,!1)){var r,i,a=e.nodeType===o;if(ve(n,e,!0))if(a)e===c&&u!==e.length&&f.push([t,e.splitText(u)]),e===d&&l&&(e.splitText(l),f.push([t,e]));else for(r=e.firstChild;r;r=i)i=r.nextSibling,p(r,t);else"INPUT"===e.nodeName||a&&!e.data||f.push([t,e])}},g=Array.prototype.filter.call(s.getElementsByTagName(e),function(o){return ve(n,o,!0)&&Y(o,e,t)});return r||g.forEach(function(e){p(e,e)}),f.forEach(function(e){var t=e[0].cloneNode(!1),n=e[1];ne(n,t),t.appendChild(n)}),g.forEach(function(e){ne(e,oe(e))}),this._getRangeAndRemoveBookmark(n),i&&n.collapse(!1),de(s,n),n},Je.changeFormat=function(e,t,n,o){return n||(n=this.getSelection())?(this.saveUndoState(n),t&&(n=this._removeFormat(t.tag.toUpperCase(),t.attributes||{},n,o)),e&&(n=this._addFormat(e.tag.toUpperCase(),e.attributes||{},n)),this.setSelection(n),this._updatePath(n,!0),x||this._docWasChanged(),this):this};var lt={DT:"DD",DD:"DT",LI:"LI",PRE:"PRE"},ct=function(e,t,n,o){var r=lt[t.nodeName],i=null,a=se(n,o,t.parentNode,e._root),s=e._config;return r||(r=s.blockTag,i=s.blockAttributes),Y(a,r,i)||(t=re(a.ownerDocument,r,i),a.dir&&(t.dir=a.dir),ne(a,t),t.appendChild(oe(a)),a=t),a};Je.forEachBlock=function(e,t,n){if(!n&&!(n=this.getSelection()))return this;t&&this.saveUndoState(n);var o=this._root,r=Ce(n,o),i=Se(n,o);if(r&&i)do{if(e(r)||r===i)break}while(r=$(r,o));return t&&(this.setSelection(n),this._updatePath(n,!0),x||this._docWasChanged()),this},Je.modifyBlocks=function(e,t){if(!t&&!(t=this.getSelection()))return this;this._recordUndoState(t,this._isInUndoState);var n,o=this._root;return be(t,o),Ne(t,o,o,o),n=pe(t,o,o),fe(t,e.call(this,n)),t.endOffset<t.endContainer.childNodes.length&&ce(t.endContainer.childNodes[t.endOffset],o),ce(t.startContainer.childNodes[t.startOffset],o),this._getRangeAndRemoveBookmark(t),this.setSelection(t),this._updatePath(t,!0),x||this._docWasChanged(),this};var ht=function(e){var t=this._root,n=e.querySelectorAll("blockquote");return Array.prototype.filter.call(n,function(e){return!X(e.parentNode,t,"BLOCKQUOTE")}).forEach(function(e){ne(e,oe(e))}),e},ut=function(){return this.createDefaultBlock([this.createElement("INPUT",{id:st,type:"hidden"}),this.createElement("INPUT",{id:dt,type:"hidden"})])},ft=function(e,t,n){for(var o,r,i,a,s=Z(t,e._root),d=e._config.tagAttributes,l=d[n.toLowerCase()],c=d.li;o=s.nextNode();)"LI"===o.parentNode.nodeName&&(o=o.parentNode,s.currentNode=o.lastChild),"LI"!==o.nodeName?(a=e.createElement("LI",c),o.dir&&(a.dir=o.dir),(i=o.previousSibling)&&i.nodeName===n?(i.appendChild(a),te(o)):ne(o,e.createElement(n,l,[a])),a.appendChild(oe(o)),s.currentNode=a):(r=(o=o.parentNode).nodeName)!==n&&/^[OU]L$/.test(r)&&ne(o,e.createElement(n,l,[oe(o)]))},pt=function(e,t){for(var n=e.commonAncestorContainer,o=e.startContainer,r=e.endContainer;n&&n!==t&&!/^[OU]L$/.test(n.nodeName);)n=n.parentNode;if(!n||n===t)return null;for(o===n&&(o=o.childNodes[e.startOffset]),r===n&&(r=r.childNodes[e.endOffset]);o&&o.parentNode!==n;)o=o.parentNode;for(;r&&r.parentNode!==n;)r=r.parentNode;return[n,o,r]};Je.increaseListLevel=function(e){if(!e&&!(e=this.getSelection()))return this.focus();var t=this._root,n=pt(e,t);if(!n)return this.focus();var o=n[0],r=n[1],i=n[2];if(!r||r===o.firstChild)return this.focus();this._recordUndoState(e,this._isInUndoState);var a,s,d=o.nodeName,l=r.previousSibling;l.nodeName!==d&&(a=this._config.tagAttributes[d.toLowerCase()],l=this.createElement(d,a),o.insertBefore(l,r));do{s=r===i?null:r.nextSibling,l.appendChild(r)}while(r=s);return(s=l.nextSibling)&&ce(s,t),this._getRangeAndRemoveBookmark(e),this.setSelection(e),this._updatePath(e,!0),x||this._docWasChanged(),this.focus()},Je.decreaseListLevel=function(e){if(!e&&!(e=this.getSelection()))return this.focus();var t=this._root,n=pt(e,t);if(!n)return this.focus();var o=n[0],r=n[1],i=n[2];r||(r=o.firstChild),i||(i=o.lastChild),this._recordUndoState(e,this._isInUndoState);var a,s=o.parentNode,d=i.nextSibling?se(o,i.nextSibling,s,t):o.nextSibling;if(s!==t&&"LI"===s.nodeName){for(s=s.parentNode;d;)a=d.nextSibling,i.appendChild(d),d=a;d=o.parentNode.nextSibling}var l=!/^[OU]L$/.test(s.nodeName);do{a=r===i?null:r.nextSibling,o.removeChild(r),l&&"LI"===r.nodeName&&(r=this.createDefaultBlock([oe(r)])),s.insertBefore(r,d)}while(r=a);return o.firstChild||te(o),d&&ce(d,t),this._getRangeAndRemoveBookmark(e),this.setSelection(e),this._updatePath(e,!0),x||this._docWasChanged(),this.focus()},Je._ensureBottomLine=function(){var e=this._root,t=e.lastElementChild;t&&t.nodeName===this._config.blockTag&&K(t)||e.appendChild(this.createDefaultBlock())},Je.setKeyHandler=function(e,t){return this._keyHandlers[e]=t,this},Je._getHTML=function(){return this._root.innerHTML},Je._setHTML=function(e){var t=this._root,n=t;n.innerHTML=e;do{ie(n,t)}while(n=$(n,t));this._ignoreChange=!0},Je.getHTML=function(e){var t,n,o,r,i,a,s=[];if(e&&(a=this.getSelection())&&this._saveRangeToBookmark(a),b)for(n=t=this._root;n=$(n,t);)n.textContent||n.querySelector("BR")||(o=this.createElement("BR"),n.appendChild(o),s.push(o));if(r=this._getHTML().replace(/\u200B/g,""),b)for(i=s.length;i--;)te(s[i]);return a&&this._getRangeAndRemoveBookmark(a),r},Je.setHTML=function(e){var t,n,o,r=this._config,i=r.isSetHTMLSanitized?r.sanitizeToDOMFragment:null,a=this._root;"function"==typeof i?n=i(e,!1,this):((t=this.createElement("DIV")).innerHTML=e,(n=this._doc.createDocumentFragment()).appendChild(oe(t))),He(n),Ke(n,a,!1),ae(n,a);for(var s=n;s=$(s,a);)ie(s,a);for(this._ignoreChange=!0;o=a.lastChild;)a.removeChild(o);a.appendChild(n),ie(a,a),this._undoIndex=-1,this._undoStack.length=0,this._undoStackLength=0,this._isInUndoState=!1;var d=this._getRangeAndRemoveBookmark()||this._createRange(a.firstChild,0);return this.saveUndoState(d),this._lastSelection=d,ot.call(this),this._updatePath(d,!0),this},Je.insertElement=function(e,t){if(t||(t=this.getSelection()),t.collapse(!0),q(e))fe(t,e),t.setStartAfter(e);else{for(var n,o=this._root,r=Ce(t,o)||o;r!==o&&!r.nextSibling;)r=r.parentNode;r!==o&&(n=se(r.parentNode,r.nextSibling,o,o)),n?o.insertBefore(e,n):(o.appendChild(e),n=this.createDefaultBlock(),o.appendChild(n)),t.setStart(n,0),t.setEnd(n,0),_e(t)}return this.focus(),this.setSelection(t),this._updatePath(t),x||this._docWasChanged(),this},Je.insertImage=function(e,t){var n=this.createElement("IMG",Ye({src:e},t,!0));return this.insertElement(n),n},Je.insertAudio=function(e,t){var n=this.createElement("audio",Ye({src:e,controls:"controls"},t,!0));return this.insertElement(n),n},Je.insertVideo=function(e,t){var n=this.createElement("video",Ye({src:e,controls:"controls"},t,!0));return this.insertElement(n),n},Je.insertFile=function(e,t){var n=e.substr(e.lastIndexOf("/")+1),o=this.createElement("a",Ye({href:e,target:"_blank",download:n},t,!0));return o.innerText=n,this.insertElement(o),o};var gt=/\b((?:(?:ht|f)tps?:\/\/|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,}\/)(?:[^\s()<>]+|\([^\s()<>]+\))+(?:\((?:[^\s()<>]+|(?:\([^\s()<>]+\)))*\)|[^\s!()\[\]{};:'".,<>?«»“”‘’]))|([\w\-.%+]+@(?:[\w\-]+\.)+[A-Z]{2,}\b)/i,mt=function(e,t,n){for(var o,r,i,a,s,d,l,c=e.ownerDocument,h=new D(e,4,function(e){return!X(e,t,"A")},!1),u=n._config.tagAttributes.a;o=h.nextNode();)for(r=o.data,i=o.parentNode;a=gt.exec(r);)d=(s=a.index)+a[0].length,s&&(l=c.createTextNode(r.slice(0,s)),i.insertBefore(l,o)),(l=n.createElement("A",Ye({href:a[1]?/^(?:ht|f)tps?:/.test(a[1])?a[1]:"http://"+a[1]:"mailto:"+a[2]},u,!1))).textContent=r.slice(s,d),i.insertBefore(l,o),o.data=r=r.slice(d)};Je.insertHTML=function(e,t){var n,o,r,i,a,s,d,l=this._config,c=l.isInsertedHTMLSanitized?l.sanitizeToDOMFragment:null,h=this.getSelection(),u=this._doc;"function"==typeof c?i=c(e,t,this):(t&&(n=e.indexOf("\x3c!--StartFragment--\x3e"),o=e.lastIndexOf("\x3c!--EndFragment--\x3e"),n>-1&&o>-1&&(e=e.slice(n+20,o))),/<\/td>((?!<\/tr>)[\s\S])*$/i.test(e)&&(e="<TR>"+e+"</TR>"),/<\/tr>((?!<\/table>)[\s\S])*$/i.test(e)&&(e="<TABLE>"+e+"</TABLE>"),(r=this.createElement("DIV")).innerHTML=e,(i=u.createDocumentFragment()).appendChild(oe(r))),this.saveUndoState(h);try{for(a=this._root,s=i,d={fragment:i,preventDefault:function(){this.defaultPrevented=!0},defaultPrevented:!1},mt(i,i,this),He(i),Ke(i,a,!1),We(i),i.normalize();s=$(s,i);)ie(s,a);t&&this.fireEvent("willPaste",d),d.defaultPrevented||(me(h,d.fragment,a),x||this._docWasChanged(),h.collapse(!1),this._ensureBottomLine()),this.setSelection(h),this._updatePath(h,!0),t&&this.focus()}catch(e){this.didError(e)}return this};var vt=function(e){return e.split("&").join("&amp;").split("<").join("&lt;").split(">").join("&gt;").split('"').join("&quot;")};Je.insertPlainText=function(e,t){var n,o,r,i,a=e.split("\n"),s=this._config,d=s.blockTag,l=s.blockAttributes,c="</"+d+">",h="<"+d;for(n in l)h+=" "+n+'="'+vt(l[n])+'"';for(h+=">",o=0,r=a.length;o<r;o+=1)i=a[o],i=vt(i).replace(/ (?= )/g,"&nbsp;"),a[o]=h+(i||"<BR>")+c;return this.insertHTML(a.join(""),t)};var _t=function(e,t,n){return function(){return this[e](t,n),this.focus()}};Je.addStyles=function(e){if(e){var t=this._doc.documentElement.firstChild,n=this.createElement("STYLE",{type:"text/css"});n.appendChild(this._doc.createTextNode(e)),t.appendChild(n)}return this},Je.bold=_t("changeFormat",{tag:"B"}),Je.italic=_t("changeFormat",{tag:"I"}),Je.underline=_t("changeFormat",{tag:"U"}),Je.strikethrough=_t("changeFormat",{tag:"S"}),Je.subscript=_t("changeFormat",{tag:"SUB"},{tag:"SUP"}),Je.superscript=_t("changeFormat",{tag:"SUP"},{tag:"SUB"}),Je.removeBold=_t("changeFormat",null,{tag:"B"}),Je.removeItalic=_t("changeFormat",null,{tag:"I"}),Je.removeUnderline=_t("changeFormat",null,{tag:"U"}),Je.removeStrikethrough=_t("changeFormat",null,{tag:"S"}),Je.removeSubscript=_t("changeFormat",null,{tag:"SUB"}),Je.removeSuperscript=_t("changeFormat",null,{tag:"SUP"}),Je.makeLink=function(e,t){var n=this.getSelection();if(n.collapsed){var o=e.indexOf(":")+1;if(o)for(;"/"===e[o];)o+=1;fe(n,this._doc.createTextNode(e.slice(o)))}return t=Ye(Ye({href:e},t,!0),this._config.tagAttributes.a,!1),this.changeFormat({tag:"A",attributes:t},{tag:"A"},n),this.focus()},Je.removeLink=function(){return this.changeFormat(null,{tag:"A"},this.getSelection(),!0),this.focus()},Je.setFontFace=function(e){return this.changeFormat(e?{tag:"SPAN",attributes:{class:l,style:"font-family: "+e+", sans-serif;"}}:null,{tag:"SPAN",attributes:{class:l}}),this.focus()},Je.setFontSize=function(e){return this.changeFormat(e?{tag:"SPAN",attributes:{class:c,style:"font-size: "+("number"==typeof e?e+"px":e)}}:null,{tag:"SPAN",attributes:{class:c}}),this.focus()},Je.setTextColour=function(e){return this.changeFormat(e?{tag:"SPAN",attributes:{class:d,style:"color:"+e}}:null,{tag:"SPAN",attributes:{class:d}}),this.focus()},Je.setHighlightColour=function(e){return this.changeFormat(e?{tag:"SPAN",attributes:{class:s,style:"background-color:"+e}}:e,{tag:"SPAN",attributes:{class:s}}),this.focus()},Je.setTextAlignment=function(e){return this.forEachBlock(function(t){var n=t.className.split(/\s+/).filter(function(e){return!!e&&!/^align/.test(e)}).join(" ");e?(t.className=n+" align-"+e,t.style.textAlign=e):(t.className=n,t.style.textAlign="")},!0),this.focus()},Je.setTextDirection=function(e){return this.forEachBlock(function(t){e?t.dir=e:t.removeAttribute("dir")},!0),this.focus()},Je.removeAllFormatting=function(e){if(!e&&!(e=this.getSelection())||e.collapsed)return this;for(var t=this._root,n=e.commonAncestorContainer;n&&!K(n);)n=n.parentNode;if(n||(be(e,t),n=t),n.nodeType===o)return this;this.saveUndoState(e),Ne(e,n,n,t);for(var r,i,a=n.ownerDocument,s=e.startContainer,d=e.startOffset,l=e.endContainer,c=e.endOffset,h=a.createDocumentFragment(),u=a.createDocumentFragment(),f=se(l,c,n,t),p=se(s,d,n,t);p!==f;)r=p.nextSibling,h.appendChild(p),p=r;return function e(t,n,r){var i,a;for(i=n.firstChild;i;i=a){if(a=i.nextSibling,q(i)){if(i.nodeType===o||"BR"===i.nodeName||"IMG"===i.nodeName){r.appendChild(i);continue}}else if(K(i)){r.appendChild(t.createDefaultBlock([e(t,i,t._doc.createDocumentFragment())]));continue}e(t,i,r)}return r}(this,h,u),u.normalize(),p=u.firstChild,r=u.lastChild,i=n.childNodes,p?(n.insertBefore(u,f),d=B.call(i,p),c=B.call(i,r)+1):c=d=B.call(i,f),e.setStart(n,d),e.setEnd(n,c),de(n,e),_e(e),this.setSelection(e),this._updatePath(e,!0),this.focus()},Je.increaseQuoteLevel=_t("modifyBlocks",function(e){return this.createElement("BLOCKQUOTE",this._config.tagAttributes.blockquote,[e])}),Je.decreaseQuoteLevel=_t("modifyBlocks",ht),Je.makeUnorderedList=_t("modifyBlocks",function(e){return ft(this,e,"UL"),e}),Je.makeOrderedList=_t("modifyBlocks",function(e){return ft(this,e,"OL"),e}),Je.removeList=_t("modifyBlocks",function(e){var t,n,o,r,i,a=e.querySelectorAll("UL, OL"),s=e.querySelectorAll("LI"),d=this._root;for(t=0,n=a.length;t<n;t+=1)ae(r=oe(o=a[t]),d),ne(o,r);for(t=0,n=s.length;t<n;t+=1)K(i=s[t])?ne(i,this.createDefaultBlock([oe(i)])):(ae(i,d),ne(i,oe(i)));return e}),Xe.isInline=q,Xe.isBlock=K,Xe.isContainer=G,Xe.getBlockWalker=Z,Xe.getPreviousBlock=j,Xe.getNextBlock=$,Xe.areAlike=V,Xe.hasTagAttributes=Y,Xe.getNearest=X,Xe.isOrContains=J,Xe.detach=te,Xe.replaceWith=ne,Xe.empty=oe,Xe.getNodeBefore=he,Xe.getNodeAfter=ue,Xe.insertNodeInRange=fe,Xe.extractContentsOfRange=pe,Xe.deleteContentsOfRange=ge,Xe.insertTreeFragmentIntoRange=me,Xe.isNodeContainedInRange=ve,Xe.moveRangeBoundariesDownTree=_e,Xe.moveRangeBoundariesUpTree=Ne,Xe.getStartBlockOfRange=Ce,Xe.getEndBlockOfRange=Se,Xe.contentWalker=ye,Xe.rangeDoesStartAtBlockBoundary=Te,Xe.rangeDoesEndAtBlockBoundary=Ee,Xe.expandRangeToBlockBoundaries=be,Xe.onPaste=Qe,Xe.addLinks=mt,Xe.splitBlock=ct,Xe.startSelectionId=st,Xe.endSelectionId=dt,"object"==typeof exports?module.exports=Xe:"function"==typeof define&&define.amd?define(function(){return Xe}):(u.Squire=Xe,top!==u&&"true"===e.documentElement.getAttribute("data-squireinit")&&(u.editor=new Xe(e),u.onEditorLoad&&(u.onEditorLoad(u.editor),u.onEditorLoad=null)))}(document);`
		jsfiletotal += "<script src='/squire.js' ></script>'"
		jstotal += "\nvar editarea={};function HsbToRgb(hsb){hsb[0]=parseFloat(hsb[0]);hsb[1]=parseFloat(hsb[1]);hsb[2]=parseFloat(hsb[2]);rgb=Array();offset=240;j=0;for(;j<3;){x=parseFloat(Math.abs(parseFloat(parseInt(hsb[0]+parseFloat(offset))%360-240)));if(x <= 60){rgb[j]=255;}else if(60<x&&x<120){rgb[j]=((1-(x-60)/60)*255);}else{rgb[j]=0;}j++;offset-= 120;}for(j=0;j<3;j++){rgb[j]+=(255-rgb[j])*(1-hsb[1]);}for(j=0;j<3;j++){rgb[j]*=hsb[2];}return rgb;}function hex2(val){vals=parseInt(val).toString(16);if(vals.length<2)vals='0'+vals;return vals;}Squire.prototype.removeHeader = function() {return this.modifyBlocks( function( frag ) {var output = this._doc.createDocumentFragment();var block = frag;while ( block = Squire.getNextBlock( block ) ) {output.appendChild(this.createElement( 'p', {'class':'paragraph'},[ Squire.empty( block )] ));}return output;});};Squire.prototype.makeHeader = function() {return this.modifyBlocks( function( frag ) {var output = this._doc.createDocumentFragment();var block = frag;while ( block = Squire.getNextBlock( block ) ) {output.appendChild(this.createElement( 'h3', [ Squire.empty( block ) ] ));}return output;});};"

		/*
					fontdlg := NewDialogBox("fontdlg", "", nil, T(T_FONT, countryid), "", true, 2)
					fontdlg.Add(NewInput("fontdlgeaid", "", nil, "hidden", "", "", "", ""))
					hbox := NewHBox("", "", nil, "")
					hbox.Add(NewLabel("", "", nil, T(T_FONTSIZE, countryid), "left"))
					hbox.Add(NewInput("editorfontsize", "", nil, "", "", T(T_EXAMPLE, countryid)+":10pt", "", ""))
					fontdlg.Add(hbox)
					hbox = NewHBox("", "", nil, "")
					hbox.Add(NewLabel("", "", nil, T(T_FONTCOLOR, countryid), "left"))
					hbox.Add(NewInput("editorfontcolor' style='color:red;", "", nil, "", "", T(T_FORMAT, countryid)+"hsb:8,1,0.9;rgb:#RRGGBB", "", ""))
					fontdlg.Add(hbox)
					hbox = NewHBox("", "", nil, "")
					hbox.Add(NewLabel("", "", nil, T(T_BACKGROUNDCOLOR, countryid), "left"))
					hbox.Add(NewInput("editorbackgroundcolor' style='background-color:#FFFF00;", "", nil, "", "", T(T_FORMAT, countryid)+"hsb:99,0.8,1;rgb:#RRGGBB", "", ""))
					fontdlg.Add(hbox)
					fontdlg.Add(NewButton("", "", nil, T(T_OK, countryid), `
					editarea[d('fontdlgeaid').value]['setFontSize'](d('editorfontsize').value);
			value = d('editorfontcolor').value;
			if(value.indexOf(',')!=value.lastIndexOf(',')){
				hsb=value.split(",");
				rgb=HsbToRgb(hsb);
				value='#'+hex2(rgb[0])+hex2(rgb[1])+hex2(rgb[2]);
			}
					editarea[d('fontdlgeaid').value]['setTextColour'](value);
			value = d('editorbackgroundcolor').value;
			if(value.indexOf(',')!=value.lastIndexOf(',')){
				hsb=value.split(",");
				rgb=HsbToRgb(hsb);
				value='#'+hex2(rgb[0])+hex2(rgb[1])+hex2(rgb[2]);
			}
					editarea[d('fontdlgeaid').value]['setHighlightColour'](value);
					editarea[d('fontdlgeaid').value].focus();
					document.location.href="#fontdlg_close";
					d('fontdlg_openModal').style.display='none';
					`))
					fontdlghtml := fontdlg.ToHtml(countryid)
					fmt.Println(fontdlghtml)
					fontdlghtml = fontdlghtml[strings.Index(fontdlghtml, "</div>")+6:]
					html += fontdlghtml
		*/

		/*
					insertfiledlg := NewDialogBox("insertfiledlg", "", nil, T(T_INSERTFILE, countryid), "", true, 2)
					insertfiledlg.Add(NewInput("insertfiledlgeaid", "", nil, "hidden", "", "", "", ""))
					insertfiledlg.Add(NewInput("insertfileurl", "", nil, "hidden", "", "", "", ""))
					hbox = NewHBox("", "", nil, "")
					//hbox.Add(NewButton("", "", nil, T(T_IMAGEFILE, countryid), "d('insertfileedit').click();"))
					hbox.Add(NewInput("insertfileedit", "", nil, "file", "multiple", "", "", ""))
					hbox.Add(NewLabel("insertfiledlguploadpt' style='color:red;", "", nil, "", "center"))
					insertfiledlg.Add(hbox)
					hbox = NewHBox("", "", nil, "")
					hbox.Add(NewLabel("", "", nil, T(T_KEYWORDS, countryid), "left"))
					hbox.Add(NewInput("insertfilekeywords", "", nil, "", "", T(T_SEPARATORBYCOMMA, countryid), "", ""))
					insertfiledlg.Add(hbox)
					insertfiledlg.Add(NewButton("", "", nil, T(T_UPLOAD, countryid), `
			var form = new FormData();
			form.append("keywords", d('insertfilekeywords').value);
			form.append("imagefile", d("insertfileedit").files[0]);
			postfile(d('insertfileurl').value,form,function(ret){
				d('insertfiledlguploadpt').innerText="100%";
				retjn=JSON.parse(ret);
				if(retjn.length>=3){
					d('insertimagelink').value=retjn[2];
				}
			},function(progress){
				if(progress>=0 && progress<=1){
					d('insertfiledlguploadpt').innerText=(progress*100)+"%";
				}else{
					d('insertfiledlguploadpt').innerText=progress;
				}
			});
					`))
					hbox = NewHBox("", "", nil, "")
					hbox.Add(NewLabel("", "", nil, T(T_LINK, countryid), "left"))
					hbox.Add(NewInput("insertimagelink", "", nil, "", "", T(T_SEPARATORBYCOMMA, countryid), "", ""))
					insertfiledlg.Add(hbox)
					insertfiledlg.Add(NewButton("", "", nil, T(T_OK, countryid), `
			urlls=d('insertimagelink').value.split(",");
			for(j=0;j<urlls.length;j++){
				if(urlls[j].trim()=="")continue;
				var ext=urlls[j].substr(urlls[j].lastIndexOf(".")+1).toLowerCase()
				if(ext=="jpg" || ext=="jpeg" || ext=="png" || ext=="gif" || ext=="svg" || ext=="bmp" || ext=="ico"){
				editarea[d('insertfiledlgeaid').value]['insertImage'](urlls[j],{"title":d('insertfilekeywords').value,"width":"100%","onload":"if(this.parentNode.parentNode.offsetWidth-15>=this.naturalWidth){this.style.width='auto';}else{this.style.width='100%';}"});
				}else if(ext=="mp3" || ext=="ogg" || ext=="acc" || ext=="m4a" || ext=="wav" || ext=="midi"){
				editarea[d('insertfiledlgeaid').value]['insertAudio'](urlls[j],{"title":d('insertfilekeywords').value});
				}else if(ext=="mp4" || ext=="ogg" || ext=="acc" || ext=="m4a" || ext=="wav" || ext=="midi"){
				editarea[d('insertfiledlgeaid').value]['insertVideo'](urlls[j],{"title":d('insertfilekeywords').value,"oncanplay":"if(this.parentNode.offsetWidth-15>=this.videoWidth){this.style.width='auto';}else{this.style.width='100%';}"});
				}else{
				editarea[d('insertfiledlgeaid').value]['insertFile'](urlls[j],{"title":d('insertfilekeywords').value});
				}
			}
			document.location.href="#fontdlg_close";`))
					insertimagehtml := insertfiledlg.ToHtml(countryid)
					insertimagehtml = insertimagehtml[strings.Index(insertimagehtml, "</div>")+6:]
					html += insertimagehtml
		*/

		//jsfiletotal += `<script type="text/javascript" src="/js/squiry.js"></script>`
		//jstotal
	}

	if strings.Index(html, "new DropLoad") != -1 {
		jstotal += `;var droploadls={};`
		jstotal += `(function(k,e){var i=window;var k=document;function DropLoad(q,p){this.element=q;this.upInsertDOM=false;this.loading=false;this.isLockUp=false;this.isLockDown=false;this.isData=true;this._scrollTop=0;this._threshold=0;this.init(p)}function a(){var p=document.documentElement.scrollTop;if(p){return p}else{return document.body.scrollTop}}DropLoad.prototype.init=function(q){var r=this;if(q.scrollArea===e){q.scrollArea=r.element}if(q.autoLoad===e){q.autoLoad=true}if(q.distance===e){q.distance=50}if(q.threshold===e){q.threshold=""}if(q.loadUpFn===e){q.loadUpFn=""}if(q.loadDownFn===e){q.loadDownFn=""}r.opts=q;if(r.opts.loadDownFn!=""){var p=h('<div class="'+r.opts.domDown.domClass+'">'+r.opts.domDown.domRefresh+"</div>");r.element.appendChild(p);r.domDown=p}if(!!r.domDown&&r.opts.threshold===""){r._threshold=Math.floor(r.domDown.offsetHeight*1/3)}else{r._threshold=r.opts.threshold}r._threshold=50;if(r.opts.scrollArea==i){r.scrollArea=window;r._scrollContentHeight=document.body.scrollHeight;r._scrollWindowHeight=k.documentElement.clientHeight}else{r.scrollArea=r.opts.scrollArea;r._scrollContentHeight=r.element.scrollHeight;r._scrollWindowHeight=r.element.offsetHeight}f(r);r.element.ontouchstart=function(s){if(!r.loading){g(s);j(s,r)}};r.element.ontouchmove=function(s){if(!r.loading){g(s,r);l(s,r)}};r.element.ontouchend=function(){if(!r.loading){b(r)}};r.scrollArea.onscroll=function(){o(r);if(r.opts.scrollArea==i){r._scrollTop=a()}else{r._scrollTop=r.scrollArea.scrollTop}if(r.opts.loadDownFn!=""&&!r.loading&&!r.isLockDown&&(r._scrollContentHeight-r._threshold)<=(r._scrollWindowHeight+r._scrollTop+10)){n(r)}}};DropLoad.prototype.startLoad=function(){var p=this;clearTimeout(p.timer);p.timer=setTimeout(function(){if(p.opts.scrollArea==i){p._scrollWindowHeight=i.innerHeight}else{p._scrollWindowHeight=p.element.offsetHeight}f(p)},150)};function g(p){if(!p.touches){p.touches=p.originalEvent.touches}}function j(q,p){p._startY=q.touches[0].pageY;if(p.opts.scrollArea==i){p.touchScrollTop=a()}else{p.touchScrollTop=p.scrollArea.scrollTop}}function l(r,p){p._curY=r.touches[0].pageY;p._moveY=p._curY-p._startY;if(p._moveY>0){p.direction="down"}else{if(p._moveY<0){p.direction="up"}}var q=Math.abs(p._moveY);if(p.opts.loadUpFn!=""&&p.touchScrollTop<=0&&p.direction=="down"&&!p.isLockUp){r.preventDefault();p.domUp=c(p.opts.domUp.domClass);if(!p.upInsertDOM){p.element.innerHTML='<div class="'+p.opts.domUp.domClass+'"></div>'+p.element.innerHTML;p.upInsertDOM=true}m(p.domUp,0);if(q<=p.opts.distance){p._offsetY=q;p.domUp.innerHTML=p.opts.domUp.domRefresh}else{if(q>p.opts.distance&&q<=p.opts.distance*2){p._offsetY=p.opts.distance+(q-p.opts.distance)*0.5;p.domUp.innerHTML=p.opts.domUp.domUpdate}else{p._offsetY=p.opts.distance+p.opts.distance*0.5+(q-p.opts.distance*2)*0.2}}p.domUp.style.height=p._offsetY}}function b(p){var q=Math.abs(p._moveY);if(p.opts.loadUpFn!=""&&p.touchScrollTop<=0&&p.direction=="down"&&!p.isLockUp){m(p.domUp,300);if(q>p.opts.distance){p.domUp.style.height=p.domUp.children[0].offsetHeight;p.domUp.innerHTML=p.opts.domUp.domLoad;p.loading=true;p.opts.loadUpFn(p)}else{p.domUp.style.height=0;p.domUp.onwebkitTransitionEnd=p.domUp.onmozTransitionEnd=p.domUp.ontransitionend=function(){p.upInsertDOM=false;this.parentNode.removeChild(this)}}p._moveY=0}}function f(p){o(p);if(p.opts.loadDownFn!=""&&p.opts.autoLoad){if((p._scrollContentHeight-p._threshold)<=p._scrollWindowHeight){n(p)}}}function o(p){if(p.opts.scrollArea==i){p._scrollContentHeight=document.body.scrollHeight}else{p._scrollContentHeight=p.element.firstChilDropLoad.children[1].offsetHeight}}function n(p){p.direction="up";p.domDown.innerHTML=p.opts.domDown.domLoad;p.loading=true;p.opts.loadDownFn(p)}DropLoad.prototype.lock=function(q){var p=this;if(q===e){if(p.direction=="up"){p.isLockDown=true}else{if(p.direction=="down"){p.isLockUp=true}else{p.isLockUp=true;p.isLockDown=true}}}else{if(q=="up"){p.isLockUp=true}else{if(q=="down"){p.isLockDown=true;p.direction="up"}}}};DropLoad.prototype.unlock=function(){var p=this;p.isLockUp=false;p.isLockDown=false;p.direction="up"};DropLoad.prototype.noData=function(p){var q=this;if(p===e||p==true){q.isData=false}else{if(p==false){q.isData=true}}};DropLoad.prototype.resetload=function(){var p=this;if(p.direction=="down"&&p.upInsertDOM){p.domUp.style.height=0;p.domUp.onwebkitTransitionEnd=p.domUp.onmozTransitionEnd=p.domUp.ontransitionend=function(){p.loading=false;p.upInsertDOM=false;this.parentNode.removeChild(this);o(p)}}else{if(p.direction=="up"){p.loading=false;if(p.isData){p.domDown.innerHTML=p.opts.domDown.domRefresh;o(p);f(p)}else{p.domDown.innerHTML=p.opts.domDown.domNoData}}}};function m(q,p){q.css({"-webkit-transition":"all "+p+"ms","transition":"all "+p+"ms"})}if(typeof exports==="object"){module.exports=DropLoad}else{if(typeof define==="function"&&define.amd){define(function(){return DropLoad})}else{i.DropLoad=DropLoad}}}(document));`
		//jstotal += `;(function(k,e){var i=window;var k=document;function DropLoad(q,p){this.element=q;this.upInsertDOM=false;this.loading=false;this.isLockUp=false;this.isLockDown=false;this.isData=true;this._scrollTop=0;this._threshold=0;this.init(p)}function a(){var p=document.documentElement.scrollTop;if(p){return p}else{return document.body.scrollTop}}DropLoad.prototype.init=function(q){var r=this;if(q.scrollArea===e){q.scrollArea=r.element}if(q.autoLoad===e){q.autoLoad=true}if(q.distance===e){q.distance=50}if(q.threshold===e){q.threshold=""}if(q.loadUpFn===e){q.loadUpFn=""}if(q.loadDownFn===e){q.loadDownFn=""}r.opts=q;if(r.opts.loadDownFn!=""){var p=h('<div class="'+r.opts.domDown.domClass+'">'+r.opts.domDown.domRefresh+"</div>");r.element.appendChild(p);r.domDown=p}if(!!r.domDown&&r.opts.threshold===""){r._threshold=Math.floor(r.domDown.offsetHeight*1/3)}else{r._threshold=r.opts.threshold}r._threshold=50;if(r.opts.scrollArea==i){r.scrollArea=window;r._scrollContentHeight=document.body.scrollHeight;r._scrollWindowHeight=k.documentElement.clientHeight}else{r.scrollArea=r.opts.scrollArea;r._scrollContentHeight=r.element.scrollHeight;r._scrollWindowHeight=r.element.offsetHeight}f(r);r.element.ontouchstart=function(s){if(!r.loading){g(s);j(s,r)}};r.element.ontouchmove=function(s){if(!r.loading){g(s,r);l(s,r)}};r.element.ontouchend=function(){if(!r.loading){b(r)}};r.scrollArea.onscroll=function(){if(r.opts.scrollArea==i){r._scrollTop=a()}else{r._scrollTop=r.scrollArea.scrollTop}if(r.opts.loadDownFn!=""&&!r.loading&&!r.isLockDown&&(r._scrollContentHeight-r._threshold)<=(r._scrollWindowHeight+r._scrollTop+10)){n(r)}}};DropLoad.prototype.startLoad=function(){var p=this;clearTimeout(p.timer);p.timer=setTimeout(function(){if(p.opts.scrollArea==i){p._scrollWindowHeight=i.innerHeight}else{p._scrollWindowHeight=p.element.offsetHeight}f(p)},150)};function g(p){if(!p.touches){p.touches=p.originalEvent.touches}}function j(q,p){p._startY=q.touches[0].pageY;if(p.opts.scrollArea==i){p.touchScrollTop=a()}else{p.touchScrollTop=p.scrollArea.scrollTop}}function l(r,p){p._curY=r.touches[0].pageY;p._moveY=p._curY-p._startY;if(p._moveY>0){p.direction="down"}else{if(p._moveY<0){p.direction="up"}}var q=Math.abs(p._moveY);if(p.opts.loadUpFn!=""&&p.touchScrollTop<=0&&p.direction=="down"&&!p.isLockUp){r.preventDefault();p.domUp=c(p.opts.domUp.domClass);if(!p.upInsertDOM){p.element.innerHTML='<div class="'+p.opts.domUp.domClass+'"></div>'+p.element.innerHTML;p.upInsertDOM=true}m(p.domUp,0);if(q<=p.opts.distance){p._offsetY=q;p.domUp.innerHTML=p.opts.domUp.domRefresh}else{if(q>p.opts.distance&&q<=p.opts.distance*2){p._offsetY=p.opts.distance+(q-p.opts.distance)*0.5;p.domUp.innerHTML=p.opts.domUp.domUpdate}else{p._offsetY=p.opts.distance+p.opts.distance*0.5+(q-p.opts.distance*2)*0.2}}p.domUp.style.height=p._offsetY}}function b(p){var q=Math.abs(p._moveY);if(p.opts.loadUpFn!=""&&p.touchScrollTop<=0&&p.direction=="down"&&!p.isLockUp){m(p.domUp,300);if(q>p.opts.distance){p.domUp.style.height=p.domUp.children[0].offsetHeight;p.domUp.innerHTML=p.opts.domUp.domLoad;p.loading=true;p.opts.loadUpFn(p)}else{p.domUp.style.height=0;p.domUp.onwebkitTransitionEnd=p.domUp.onmozTransitionEnd=p.domUp.ontransitionend=function(){p.upInsertDOM=false;this.parentNode.removeChild(this)}}p._moveY=0}}function f(p){if(p.opts.loadDownFn!=""&&p.opts.autoLoad){if((p._scrollContentHeight-p._threshold)<=p._scrollWindowHeight){n(p)}}}function o(p){if(p.opts.scrollArea==i){p._scrollContentHeight=document.body.scrollHeight}else{p._scrollContentHeight=p.element.scrollHeight}}function n(p){p.direction="up";p.domDown.innerHTML=p.opts.domDown.domLoad;p.loading=true;p.opts.loadDownFn(p)}DropLoad.prototype.lock=function(q){var p=this;if(q===e){if(p.direction=="up"){p.isLockDown=true}else{if(p.direction=="down"){p.isLockUp=true}else{p.isLockUp=true;p.isLockDown=true}}}else{if(q=="up"){p.isLockUp=true}else{if(q=="down"){p.isLockDown=true;p.direction="up"}}}};DropLoad.prototype.unlock=function(){var p=this;p.isLockUp=false;p.isLockDown=false;p.direction="up"};DropLoad.prototype.noData=function(p){var q=this;if(p===e||p==true){q.isData=false}else{if(p==false){q.isData=true}}};DropLoad.prototype.resetload=function(){var p=this;if(p.direction=="down"&&p.upInsertDOM){p.domUp.style.height=0;p.domUp.onwebkitTransitionEnd=p.domUp.onmozTransitionEnd=p.domUp.ontransitionend=function(){p.loading=false;p.upInsertDOM=false;this.parentNode.removeChild(this);o(p)}}else{if(p.direction=="up"){p.loading=false;if(p.isData){p.domDown.innerHTML=p.opts.domDown.domRefresh;o(p);f(p)}else{p.domDown.innerHTML=p.opts.domDown.domNoData}}}};function m(q,p){q.css({"-webkit-transition":"all "+p+"ms","transition":"all "+p+"ms"})}if(typeof exports==="object"){module.exports=DropLoad}else{if(typeof define==="function"&&define.amd){define(function(){return DropLoad})}else{i.DropLoad=DropLoad}}}(document));`
		//jstotal += `;(function(doc,undefined){var win=window;var doc=document;function DropLoad(element,options){this.element=element;this.upInsertDOM=false;this.loading=false;this.isLockUp=false;this.isLockDown=false;this.isData=true;this._scrollTop=0;this._threshold=0;this.init(options)}function wscrollTop(){var t=document.documentElement.scrollTop;if(t){return t}else{return document.body.scrollTop}}DropLoad.prototype.init=function(options){var me=this;if(options.scrollArea===undefined){options.scrollArea=me.element}if(options.autoLoad===undefined){options.autoLoad=true}if(options.distance===undefined){options.distance=50}if(options.threshold===undefined){options.threshold=""}if(options.loadUpFn===undefined){options.loadUpFn=""}if(options.loadDownFn===undefined){options.loadDownFn=""}me.opts=options;if(me.opts.loadDownFn!=""){var downnode=h('<div class="'+me.opts.domDown.domClass+'">'+me.opts.domDown.domRefresh+"</div>");me.element.appendChild(downnode);me.domDown=downnode}if(!!me.domDown&&me.opts.threshold===""){me._threshold=Math.floor(me.domDown.offsetHeight*1/3)}else{me._threshold=me.opts.threshold}me._threshold=50;if(me.opts.scrollArea==win){me.scrollArea=window;me._scrollContentHeight=document.body.scrollHeight;me._scrollWindowHeight=doc.documentElement.clientHeight}else{me.scrollArea=me.opts.scrollArea;me._scrollContentHeight=me.element.scrollHeight;me._scrollWindowHeight=me.element.offsetHeight}fnAutoLoad(me);me.element.ontouchstart=function(e){if(!me.loading){fnTouches(e);fnTouchstart(e,me)}};me.element.ontouchmove=function(e){if(!me.loading){fnTouches(e,me);fnTouchmove(e,me)}};me.element.ontouchend=function(){if(!me.loading){fnTouchend(me)}};me.scrollArea.onscroll=function(){if(me.opts.scrollArea==win){me._scrollTop=wscrollTop()}else{me._scrollTop=me.scrollArea.scrollTop}if(me.opts.loadDownFn!=""&&!me.loading&&!me.isLockDown&&(me._scrollContentHeight-me._threshold)<=(me._scrollWindowHeight+me._scrollTop+10)){loadDown(me)}}};DropLoad.prototype.startLoad=function(){var me=this;clearTimeout(me.timer);me.timer=setTimeout(function(){if(me.opts.scrollArea==win){me._scrollWindowHeight=win.innerHeight}else{me._scrollWindowHeight=me.element.offsetHeight}fnAutoLoad(me)},150)};function fnTouches(e){if(!e.touches){e.touches=e.originalEvent.touches}}function fnTouchstart(e,me){me._startY=e.touches[0].pageY;if(me.opts.scrollArea==win){me.touchScrollTop=wscrollTop()}else{me.touchScrollTop=me.scrollArea.scrollTop}}function fnTouchmove(e,me){me._curY=e.touches[0].pageY;me._moveY=me._curY-me._startY;if(me._moveY>0){me.direction="down"}else{if(me._moveY<0){me.direction="up"}}var _absMoveY=Math.abs(me._moveY);if(me.opts.loadUpFn!=""&&me.touchScrollTop<=0&&me.direction=="down"&&!me.isLockUp){e.preventDefault();me.domUp=c(me.opts.domUp.domClass);if(!me.upInsertDOM){me.element.innerHTML='<div class="'+me.opts.domUp.domClass+'"></div>'+me.element.innerHTML;me.upInsertDOM=true}fnTransition(me.domUp,0);if(_absMoveY<=me.opts.distance){me._offsetY=_absMoveY;me.domUp.innerHTML=me.opts.domUp.domRefresh}else{if(_absMoveY>me.opts.distance&&_absMoveY<=me.opts.distance*2){me._offsetY=me.opts.distance+(_absMoveY-me.opts.distance)*0.5;me.domUp.innerHTML=me.opts.domUp.domUpdate}else{me._offsetY=me.opts.distance+me.opts.distance*0.5+(_absMoveY-me.opts.distance*2)*0.2}}me.domUp.style.height=me._offsetY}}function fnTouchend(me){var _absMoveY=Math.abs(me._moveY);if(me.opts.loadUpFn!=""&&me.touchScrollTop<=0&&me.direction=="down"&&!me.isLockUp){fnTransition(me.domUp,300);if(_absMoveY>me.opts.distance){me.domUp.style.height=me.domUp.children[0].offsetHeight;me.domUp.innerHTML=me.opts.domUp.domLoad;me.loading=true;me.opts.loadUpFn(me)}else{me.domUp.style.height=0;me.domUp.onwebkitTransitionEnd=me.domUp.onmozTransitionEnd=me.domUp.ontransitionend=function(){me.upInsertDOM=false;this.parentNode.removeChild(this)}}me._moveY=0}}function fnAutoLoad(me){if(me.opts.loadDownFn!=""&&me.opts.autoLoad){if((me._scrollContentHeight-me._threshold)<=me._scrollWindowHeight){loadDown(me)}}}function fnRecoverContentHeight(me){if(me.opts.scrollArea==win){me._scrollContentHeight=document.body.scrollHeight}else{me._scrollContentHeight=me.element.scrollHeight}}function loadDown(me){me.direction="up";me.domDown.innerHTML=me.opts.domDown.domLoad;me.loading=true;me.opts.loadDownFn(me)}DropLoad.prototype.lock=function(direction){var me=this;if(direction===undefined){if(me.direction=="up"){me.isLockDown=true}else{if(me.direction=="down"){me.isLockUp=true}else{me.isLockUp=true;me.isLockDown=true}}}else{if(direction=="up"){me.isLockUp=true}else{if(direction=="down"){me.isLockDown=true;me.direction="up"}}}};DropLoad.prototype.unlock=function(){var me=this;me.isLockUp=false;me.isLockDown=false;me.direction="up"};DropLoad.prototype.noData=function(flag){var me=this;if(flag===undefined||flag==true){me.isData=false}else{if(flag==false){me.isData=true}}};DropLoad.prototype.resetload=function(){var me=this;if(me.direction=="down"&&me.upInsertDOM){me.domUp.style.height=0;me.domUp.onwebkitTransitionEnd=me.domUp.onmozTransitionEnd=me.domUp.ontransitionend=function(){me.loading=false;me.upInsertDOM=false;this.parentNode.removeChild(this);fnRecoverContentHeight(me)}}else{if(me.direction=="up"){me.loading=false;if(me.isData){me.domDown.innerHTML=me.opts.domDown.domRefresh;fnRecoverContentHeight(me);fnAutoLoad(me)}else{me.domDown.innerHTML=me.opts.domDown.domNoData}}}};function fnTransition(dom,num){dom.css({"-webkit-transition":"all "+num+"ms","transition":"all "+num+"ms"})}if(typeof exports==="object"){module.exports=DropLoad}else{if(typeof define==="function"&&define.amd){define(function(){return DropLoad})}else{win.DropLoad=DropLoad}}}(document));`
		//jsfiletotal += `<script type="text/javascript" src="/js/dropload2.js"></script>`
	}

	//if strings.Index(html, `ud.style.left="0pt";ud.style.top="0pt";ud.innerHTML="◤";`) != -1 {
	jstotal += `function getDPI(){var arrDPI=new Array;if(window.screen.deviceXDPI){arrDPI[0]=window.screen.deviceXDPI;arrDPI[1]=window.screen.deviceYDPI;}else{var tmpNode=document.createElement("DIV");tmpNode.style.cssText="width:1in;height:1in;position:absolute;left:0px;top:0px;z-index:99;visibility:hidden";document.body.appendChild(tmpNode);arrDPI[0]=parseInt(tmpNode.offsetWidth);arrDPI[1]=parseInt(tmpNode.offsetHeight);tmpNode.parentNode.removeChild(tmpNode);}return parseFloat(arrDPI);}`
	//}

	html, onloadls := toolfunc.GetAndTruncateTagContent(html, "<!--onloadbegin", "onloadend-->")
	onloadctt := ""
	for i := 0; i < len(onloadls); i++ {
		onloadctt += onloadls[i]
	}
	html, onresizels := toolfunc.GetAndTruncateTagContent(html, "<!--onresizebegin", "onresizeend-->")
	onresizectt := ""
	for i := 0; i < len(onresizels); i++ {
		onresizectt += onresizels[i]
	}
	jstotal += `
window.onload=function(){
	cls=document.getElementsByTagName('canvas');
	for(var i=0;i<cls.length;i++){
		cls[i].width=cls[i].offsetWidth;
		cls[i].height=cls[i].offsetHeight;
	}` + onloadctt + `
}
window.onresize=function(){
	cls=document.getElementsByTagName('canvas');
	for(var i=0;i<cls.length;i++){
		cls[i].width=cls[i].offsetWidth;
		cls[i].height=cls[i].offsetHeight;
	}` + onresizectt + `
}
	`
	return []byte("<!DOCTYPE html>\n<html><head><meta name='referrer' content='never'><meta name='full-screen' content='yes'><meta name='x5-fullscreen' content='true'><meta http-equiv='Content-Type' content='text/html;charset=utf-8' /><meta name='viewport' content='width=device-width, initial-scale=1.0, minimum-scale=0.5, maximum-scale=2.0, user-scalable=yes' /><title>" + title + "</title><style>" + styletotal + "</style></head>" + jsfiletotal + "<body>" + "<script type='text/javascript'>" + jstotal + "</script>" + html + "</body></html>")
}

func getClickJsExpr(ctr BaseControl, onclickjscode string) (inljscode, scriptjscode string) {
	onclickjscode = strings.Trim(onclickjscode, "\r\n\t")
	var clkjs, scriptjs string
	if strings.Index(onclickjscode, "\n") == -1 {
		clkjs = onclickjscode
	} else {
		if strings.HasPrefix(onclickjscode, "function ") {
			clkjs = strings.Trim(onclickjscode[strings.Index(onclickjscode, " ")+1:strings.Index(onclickjscode, "(")], " \r\n\t")
			if !(strings.HasPrefix(onclickjscode[strings.Index(onclickjscode, "("):], "()") || strings.HasPrefix(onclickjscode[strings.Index(onclickjscode, "("):], "(this)")) {
				clkjs = ""
			}
			scriptjs = "<script type='text/javascript'>" + onclickjscode + "</script>"
		} else {
			clkjs = "f" + fmt.Sprintf("%p", ctr)[2:]
			if strings.Index(onclickjscode, "this.") == -1 {
				scriptjs = "<script type='text/javascript'>function " + clkjs + "(){" + onclickjscode + "}</script>"
			} else {
				scriptjs = "<script type='text/javascript'>function " + clkjs + "(_this){" + onclickjscode + "}</script>"
			}
		}
		if strings.Index(onclickjscode, "this.") != -1 {
			clkjs += "(this)"
		} else {
			clkjs += "()"
		}
		if len(clkjs) > 0 && clkjs[:1] == "_" {
			clkjs = ""
		}
	}
	if clkjs == "()" {
		clkjs = ""
	}
	return clkjs, scriptjs
}

func getClickJsExprWithEvent(ctr BaseControl, onclickjscode string) (inljscode, scriptjscode string) {
	onclickjscode = strings.Trim(onclickjscode, "\r\n\t")
	var clkjs, scriptjs string
	if strings.Index(onclickjscode, "\n") == -1 {
		clkjs = onclickjscode
	} else {
		if strings.HasPrefix(onclickjscode, "function ") {
			clkjs = strings.Trim(onclickjscode[strings.Index(onclickjscode, " ")+1:strings.Index(onclickjscode, "(")], " \r\n\t")
			if !(strings.HasPrefix(onclickjscode[strings.Index(onclickjscode, "("):], "()") || strings.HasPrefix(onclickjscode[strings.Index(onclickjscode, "("):], "(this)")) {
				clkjs = ""
			}
			scriptjs = "<script type='text/javascript'>" + onclickjscode + "</script>"
		} else {
			clkjs = "f" + fmt.Sprintf("%p", ctr)[2:]
			if strings.Index(onclickjscode, "this.") == -1 {
				scriptjs = "<script type='text/javascript'>function " + clkjs + "(event){" + onclickjscode + "}</script>"
			} else {
				scriptjs = "<script type='text/javascript'>function " + clkjs + "(event,_this){" + onclickjscode + "}</script>"
			}
		}
		if strings.Index(onclickjscode, "this.") != -1 {
			clkjs += "(event,this)"
		} else {
			clkjs += "(event)"
		}
		if len(clkjs) > 0 && clkjs[:1] == "_" {
			clkjs = ""
		}
	}
	if clkjs == "()" {
		clkjs = ""
	}
	return clkjs, scriptjs
}

func percentv(val string) float64 {
	reg := regexp.MustCompile(`[0-9.]+`)
	valbyte := []byte(val)
	val3 := reg.Find(valbyte)
	val2, _ := strconv.ParseFloat(string(val3), 64)
	return val2
}

func getctrlname() string {
	return "fkljdsklaf"
}

type BoxMethod interface {
	AddBox(string)
	AddControl(string)
	ToHtml(countryid int) string
}

type ControlMethod interface {
}

type BaseControl interface {
	SetHPercent(float32)
	SetVPercent(float32)
	ToHtml(countryid int) string
	GetMinWidth() int
	GetMinHeight() int
	GetMaxWidth() int
	GetMaxHeight() int
	GetHMinPercent() float32
	GetVMinPercent() float32
	GetHMaxPercent() float32
	GetVMaxPercent() float32
	GetHExpanding() int
	GetVExpanding() int
	SetDisplay(string)
	SetFloat(string)
	CalcHeight() string
	CalcWidth() string
	SetStyle(style string) bool
}

type BaseBox struct {
	id, idname, Class, selfstyle string
	//pt unit
	Minwidth, Minheight               int
	Maxwidth, Maxheight               int
	Percentwidth, Percentheight       string
	MinPercentwidth, MinPercentheight string
	MaxPercentwidth                   string
	MaxPercentheight                  string
	//Hexpanding value:
	//1. Fixed：控件不能放大或者缩小，控件的大小就是它的sizeHint。
	//2. Minimum：控件的sizeHint为控件的最小尺寸。控件不能小于这个sizeHint，但是可以放大。
	//3. Maximum：控件的sizeHint为控件的最大尺寸，控件不能放大，但是可以缩小到它的最小的允许尺寸。
	//4. Preferred：控件的sizeHint是它的sizeHint，但是可以放大或者缩小
	//5. Expandint：控件可以自行增大或者缩小
	Hexpanding    int
	Vexpanding    int
	X, Y          int
	Alignment     int
	Margin_left   string
	Margin_right  string
	Margin_top    string
	Margin_bottom string
	child         []BaseControl
	Float         string
	Display       string
	color         []string
	attrmap       map[string]string
}

func NewBaseBox() *BaseBox {
	return &BaseBox{Percentwidth: "100%", Percentheight: "100%"}
}

func (btn *BaseBox) SetHPercent(pcwid float32) {
	btn.Percentwidth = fmt.Sprintf("%g%%", pcwid)
}

func (btn *BaseBox) SetVPercent(pchei float32) {
	btn.Percentheight = fmt.Sprintf("%g%%", pchei)
}
func (btn *BaseBox) GetMinWidth() int {
	return btn.Minwidth
}
func (btn *BaseBox) GetMinHeight() int {
	return btn.Minheight
}
func (btn *BaseBox) GetMaxWidth() int {
	return btn.Maxwidth
}
func (btn *BaseBox) GetMaxHeight() int {
	return btn.Maxheight
}

func (btn *BaseBox) GetHMinPercent() float32 {
	return float32(percentv(btn.MinPercentwidth))
}
func (btn *BaseBox) GetVMinPercent() float32 {
	return float32(percentv(btn.MinPercentheight))
}
func (btn *BaseBox) GetHMaxPercent() float32 {
	return float32(percentv(btn.MaxPercentwidth))
}
func (btn *BaseBox) GetVMaxPercent() float32 {
	return float32(percentv(btn.MaxPercentheight))
}
func (btn *BaseBox) GetHExpanding() int {
	return btn.Hexpanding
}
func (btn *BaseBox) GetVExpanding() int {
	return btn.Vexpanding
}
func (btn *BaseBox) SetDisplay(Display string) {
	btn.Display = Display
}
func (btn *BaseBox) SetFloat(Float string) {
	btn.Float = Float
}
func (btn *BaseBox) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true && btn.MaxPercentheight == "" {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}
func (btn *BaseBox) CalcWidth() string {
	if btn.Hexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxwidth), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		return strconv.FormatFloat(math.Min(percentv(btn.Percentwidth), percentv(btn.MaxPercentwidth)), 'f', 4, 32) + "%"
	}
}

func (btn *BaseBox) ToIdAttr() string {
	id := ""
	if btn.id != "" {
		if btn.id[:1] == "(" && strings.Index(btn.id, ")") != -1 {
			maxwidhei := btn.id[1:strings.Index(btn.id, ")")]
			widheils := strings.Split(maxwidhei, ",")
			if len(widheils) >= 1 && widheils[0] != "" {
				Maxwidth2, _ := strconv.ParseInt(widheils[0], 10, 32)
				btn.Maxwidth = int(Maxwidth2)
				btn.Hexpanding = 1
			}
			if len(widheils) >= 2 && widheils[1] != "" {
				Maxheight2, _ := strconv.ParseInt(widheils[1], 10, 32)
				btn.Maxheight = int(Maxheight2)
				btn.Vexpanding = 1
			}
			btn.id = btn.id[strings.Index(btn.id, ")")+1:]
		}
		if strings.Index(btn.id, "' style='") != -1 {
			btn.selfstyle = btn.id[strings.Index(btn.id, "' style='")+len("' style='"):]
			btn.id = btn.id[:strings.Index(btn.id, "' style='")]
		}
		btn.idname = btn.id
		id = " id='" + btn.id + "'"
	}
	if btn.Class != "" {
		if strings.Index(btn.Class, "' style='") != -1 {
			btn.selfstyle = btn.Class[strings.Index(btn.Class, "' style='")+len("' style='"):]
			btn.Class = btn.Class[:strings.Index(btn.Class, "' style='")]
		}
		id += " Class='" + btn.Class + "'"
	}
	return id
}

func (btn *BaseBox) SetStyle(stylestr string) bool {
	btn.selfstyle = stylestr
	return true
}

func (btn *BaseBox) SetAttr(name, value string) bool {
	if btn.attrmap == nil {
		btn.attrmap = make(map[string]string, 0)
	}
	btn.attrmap[name] = value
	return true
}

func (btn *BaseBox) AttrToText() (str string) {
	for key, val := range btn.attrmap {
		if strings.Index(val, "'") == -1 {
			str += key + "='" + val + "' "
		} else {
			str += key + "=\"" + val + "\" "
		}
	}
	if len(str) > 0 {
		str = " " + str
	}
	return str
}

type HBox struct {
	BaseBox
	onclickjs string
}

func NewHBox(id, Class string, color []string, onclickjs string) *HBox {
	return &HBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, onclickjs: strings.Trim(onclickjs, "\r\n\t ")}
}

func (btn *HBox) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true && btn.MaxPercentheight == "" {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *HBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	//set child width percent
	showcnt := 0
	for _, val := range hb.child {
		subhtml := val.ToHtml(Country_English)
		//fmt.Println("bh:", regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)), subhtml)
		if !regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)) {
			showcnt += 1
		}
	}
	onep := 100.0 / float32(showcnt)
	var ee float32
	var cnt int
	for _, val := range hb.child {
		subhtml := val.ToHtml(Country_English)
		//fmt.Println("bh:", regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)), subhtml)
		if !regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)) {
			if onep < val.GetHMaxPercent() {
				cnt += 1
			} else {
				val.SetHPercent(val.GetHMaxPercent())
				ee += onep - val.GetHMaxPercent()
			}
		}
	}
	onep += ee / float32(cnt)
	for _, val := range hb.child {
		if onep < val.GetHMaxPercent() {
			val.SetHPercent(onep)
		} else {

		}
	}
	return false
}

func (hb *HBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()

	html := "<div" + id + ` style="width:` + hb.CalcWidth() + `;height:auto;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `"`
	clkjs, scriptjs := getClickJsExpr(hb, hb.onclickjs)
	if clkjs != "" {
		html += " onclick=\"" + clkjs + "\""
	}
	html += hb.AttrToText() + `>__child__</div>`
	var childstr string
	for _, val := range hb.child {
		str := val.ToHtml(countryid)
		if len(str) > 5 && str[:5] == "<div " && !regexp.MustCompile("<div[^>]*display:inline-block;").Match([]byte(str)) {
			str = "<span " + str[4:len(str)-6] + "</span>"
		}
		childstr += str
	}
	html2 := strings.Replace(html, "__child__", childstr, -1)
	if strings.Index(hb.onclickjs, "\n") != -1 {
		html2 += scriptjs
	}
	return html2
}

type FlowBox struct {
	BaseBox
}

func NewFlowBox(id, Class string, color []string) *FlowBox {
	return &FlowBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}}
}

//FlowBox
func (hb *FlowBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	bc.SetFloat("left")
	return false
}

func (hb *FlowBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	html := "<div" + id + ` style="width:` + hb.Percentwidth + `;height:auto;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `"` + hb.AttrToText() + `>
	__child__
	</div>`
	var childstr string
	for _, val := range hb.child {
		str := val.ToHtml(countryid)
		childstr += str
	}
	return strings.Replace(html, "__child__", childstr, -1)
}

type ScrollBox struct {
	BaseBox
}

func NewScrollBox(id, Class string, color []string) *ScrollBox {
	return &ScrollBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}}
}

//ScrollBox
func (hb *ScrollBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	bc.SetFloat("left")
	return false
}
func (btn *ScrollBox) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	html := "<div" + id + ` style="width:100%;height:auto;overflow:auto;` + btn.selfstyle + `"` + btn.AttrToText() + `>
	__child__
	</div>`
	var childstr string
	for i := 0; i < len(btn.child); i++ {
		childstr += btn.child[i].ToHtml(countryid)
	}
	return strings.Replace(html, "__child__", childstr, 1)
}

type StackBox struct {
	BaseBox
	Childname []string
	Childid   []string
}

//class must unique define;that for change tab with 'tab' join class name;can append div as child;statck container is 'ctn' join class name
func NewStackBox(id, Class string, color []string) *StackBox {
	return &StackBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}}
}

//StackBox
func (hb *StackBox) Add(bc BaseControl, name, Childid string) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	hb.Childname = append(hb.Childname, name)
	hb.Childid = append(hb.Childid, Childid)
	return false
}
func (hb *StackBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	html := "<div" + id + ` style="width:` + hb.Percentwidth + `;height:` + hb.Percentheight + `;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `"` + hb.AttrToText() + `>
	__child__
	</div>`
	var childstr string
	childstr = `<div id='ctn` + hb.Class + `' Class='tab-content' style='width:100%;height:100%;'>`
	for index, val := range hb.child {
		str := val.ToHtml(countryid)
		if index == 0 {
			childstr += `<div id="` + hb.Childid[index] + `" Class='tab` + hb.Class + `' style="display:block;width:100%;height:100%;">` + str + `</div>`
		} else {
			childstr += `<div id="` + hb.Childid[index] + `" Class='tab` + hb.Class + `' style="display:none;width:100%;height:100%;">` + str + `</div>`
		}
	}
	childstr += "</div>"
	return strings.Replace(html, "__child__", childstr, 1)
}
func (hb *StackBox) HeadAsHBox() *HBox {
	hb1 := NewHBox("", "", hb.color, "")
	for index, _ := range hb.child {
		link := NewLink("", "", hb1.color, hb.Childname[index], "stackto('"+hb.Class+"',"+fmt.Sprintf("%d", index)+")", "javascript:void(0)", "", "center")
		hb1.Add(link)
	}
	hb1.MaxPercentheight = "5%"
	return hb1
}
func (hb *StackBox) HeadAsVBox() *VBox {
	vb1 := NewVBox("", "", hb.color, "")
	for index, _ := range hb.child {
		link := NewLink("", "", hb.color, hb.Childname[index], "stackto('"+hb.Class+"',"+fmt.Sprintf("%d", index)+")", "javascript:void(0)", "", "center")
		vb1.Add(link)
	}
	vb1.MaxPercentheight = "5%"
	return vb1
}

type TableBoxItem struct {
	x, y, x_expand, y_expand    int
	widthpercent, heightpercent float32
	child                       BaseControl
}

type TableBox struct {
	BaseBox
	child                    []*TableBoxItem //BaseControl
	Childminwid, Childminhei int
}

func NewTableBox(id, Class string, color []string) *TableBox {
	return &TableBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}}
}

//x ,y from 0
func (hb *TableBox) AddChild(x, x_expand, y, y_expand int, widthpercent, heightpercent float32, bc BaseControl) bool {
	hb.child = append(hb.child, &TableBoxItem{x: x, y: y, x_expand: x_expand, y_expand: y_expand, widthpercent: widthpercent, heightpercent: heightpercent, child: bc})
	return true
}

func (hb *TableBox) Delete(bc *BaseControl) bool {
	for i := 0; i < len(hb.child); i++ {
		if &hb.child[i].child == bc {
			hb.child = append(hb.child[:i], hb.child[i+1:]...)
			return true
		}
	}
	return false
}

func (hb *TableBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	html := "<table" + id + ` style="width:100%;height:auto;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `" border="1"  cellspacing="0" cellpadding="0"` + hb.AttrToText() + `>
	<tbody style="width:100%">__child__</tbody>
	</table>`
	var childstr string
	maxrow := 0
	maxcol := 0
	for i := 0; i < len(hb.child); i++ {
		if hb.child[i].x+hb.child[i].x_expand > maxrow {
			maxrow = hb.child[i].x + hb.child[i].x_expand
		}
		if hb.child[i].y+hb.child[i].y_expand > maxcol {
			maxcol = hb.child[i].y + hb.child[i].y_expand
		}
	}
	for i := 0; i < maxcol; i += 1 {
		childstr += "<tr style='width:100%;'>"
		for j := 0; j < maxrow; j += 1 {
			bfnd := false
			for k := 0; k < len(hb.child); k++ {
				if hb.child[k].x == j && hb.child[k].y == i {
					str := hb.child[k].child.ToHtml(countryid)
					//colspan rowspan
					spanstr := ""
					if hb.child[k].x_expand > 1 {
						spanstr += "rowspan=" + toolfunc.IntToStr(hb.child[k].x_expand)
					}
					if hb.child[k].x_expand > 1 {
						spanstr += "colspan=" + toolfunc.IntToStr(hb.child[k].y_expand)
					}
					if hb.child[k].widthpercent != -1 {
						spanstr += "width='" + toolfunc.Float32ToStr(hb.child[k].widthpercent) + "%'"
					}
					if hb.child[k].widthpercent != -1 {
						spanstr += "width='" + toolfunc.Float32ToStr(hb.child[k].heightpercent) + "%'"
					}
					if len(spanstr) > 0 {
						spanstr = " " + spanstr
					}
					childstr += "<td" + spanstr + ">" + str + "</td>"
					bfnd = true
				}
			}
			if bfnd == false {
				childstr += "<td></td>"
			}
		}
		childstr += "</tr>"
	}
	return strings.Replace(html, "__child__", childstr, -1)
}

type MenuBox struct {
	BaseBox
	Title string
}

func NewMenuBox(id, Class string, color []string, Title string) *MenuBox {
	return &MenuBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 1, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Title: Title}
}

//MenuBox
func (hb *MenuBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	return false
}
func (btn *MenuBox) ToHtml(countryid int) string {
	style := `<style>
.clear:after {
clear: both;
content: ".";
display: block;
height: 0;
visibility: hidden;
}
nav{
display:inline-block;
border:1px solid #505255;
border-bottom: 1px solid #282C2F;
-moz-border-radius: 5px;
-webkit-border-radius: 5px;
margin:50px;
-webkit-box-shadow:1px 1px 3px #292929;
-moz-box-shadow:1px 1px 3px #292929;
}
li{
list-style:none;
float:left;
border-right: 1px solid #2E3235;
position: relative;
/*background: -moz-linear-gradient(top, #fff, #555D5F 2% ,#555D5F 50%,#3E4245 100%);
background: -webkit-gradient(linear, 0 0, 0 100%, from(#fff), color-stop(2%, #555D5F), color-stop(50%, #555D5F),to(#3E4245));*/
background:#555D5F;
}
li:hover{
/*background: -moz-linear-gradient(top, #fff, #3E4245 2% ,#555D5F 80%,#555D5F 100%);
background: -webkit-gradient(linear, 0 0, 0 100%, from(#fff), color-stop(2%, #3E4245), color-stop(80%, #3E4245),to(#555D5F));*/
background:#3E4245;
-moz-transition: background 1s ease-out;
-webkit-transition: background 1s ease-out;
}
li a{
display:block;
height:40px;
line-height:40px;
padding:0 30px;
font-size:12px;
color:#fff;
text-shadow: 0px -1px 0px #000;
text-decoration:none;
white-space:nowrap;
border-left: 1px solid #999E9F;
border-top: 1px solid #999E9F;
-moz-border-top-left-radius: 2px;
-webkit-border-top-left-radius: 2px;
z-index:100;
}
li > a{
position:relative;
}
li.first a{
-moz-border-radius-topleft: 4px;
-moz-border-radius-bottomleft: 4px;
-webkit-border-top-left-radius: 4px;
-webkit-border-bottom-left-radius: 4px;
}
li.last{
border-right: 0 none;
}
dl{
position:absolute;
display:block;
top:40px;
left: -25px;
width:165px;
background:#222222;
-moz-border-radius: 5px;
-webkit-border-radius: 5px;
-webkit-box-shadow:1px 1px 3px #292929;
-moz-box-shadow:1px 1px 3px #292929;
z-index:10;
}
li:hover dl{
top:50px;
display:block;
width:145px;
padding:10px;
}
dl a{
background:transparent;
border:0 none;
-moz-border-radius: 5px;
-webkit-border-radius: 5px;
-moz-transition: background 0.5s ease-out; 
-webkit-transition: background 0.5s ease-out;
z-index:50;
}
dl a:hover{
color:#FFF;
background:#999E9F;
-moz-transition: background 0.3s ease-in-out, color 0.3s ease-in-out;
-webkit-transition: background 0.3s ease-in-out, color 0.3s ease-in-out;
}
dd{
margin-top:-40px;
opacity:0;
width:145px;
-webkit-transition-property:all;
/*-webkit-transition-timing-function: cubic-bezier(5,0,5,0);*/
-moz-transition-property: all;
/*-moz-transition-timing-function: cubic-bezier(5,0,5,0);*/
/*-webkit-transition-delay:5s;
-moz-transition-delay:5s;*/
}
li:hover dd{
margin-top:0;
opacity:1;
}
li dd:nth-child(1){
-webkit-transition-duration: 0.1s;
-moz-transition-duration: 0.1s;
}
li dd:nth-child(2){
-webkit-transition-duration: 0.2s;
-moz-transition-duration: 0.2s;
}
li dd:nth-child(3){
-webkit-transition-duration: 0.3s;
-moz-transition-duration: 0.3s;
}
li dd:nth-child(4){
-webkit-transition-duration: 0.4s;
-moz-transition-duration: 0.4s;
}
dt{
display:none;
margin-top:-25px;
padding-top:15px;
height:10px;
}
li:hover dt{
display:block;
}
.Darrow{
float:right;
margin:18px 10px 0 0;
border-width:5px;
border-color:#FFF transparent transparent transparent;
border-style:solid;
width:0;
height:0;
line-height:0;
overflow:hidden;
cursor:pointer;
text-shadow: 0px -1px 0px #000;
-webkit-box-shadow:0px -1px 0px #000;
-moz-box-shadow:0px -1px 0px #000;
}
.arrow{
margin:0 auto;
margin-top:-5px;
display:block;
width:10px;
height:10px;
background:#222222;
-webkit-transform: rotate(45deg);
-moz-transform: rotate(45deg);
}
</style>`
	html := `<ul Class="clear" style="display:` + btn.Display + `;` + btn.selfstyle + `" > <li><span Class="Darrow"></span> 
                    <a href="#">菜单二</a> 
                    <dl> 
                        <dt><span Class="arrow"></span></dt>
						__menuitem__
						</dl></li></ul>  `
	var childstr string
	for _, val := range btn.child {
		str := val.ToHtml(countryid)
		childstr += "<dd>" + str + "</dd>"
	}
	return style + strings.Replace(html, "__menuitem__", childstr, 1)
}

type DialogBox struct {
	BaseBox
	Title         string
	bOutSizeClose bool
	linkorbtn     int
	dialogwidth   string
}

//linkorbtn:1-left align link;2-center align link;3-right align link;4-button
func NewDialogBox(id, Class string, color []string, Title, dialogwidth string, bOutSizeClose bool, linkorbtn int) *DialogBox {
	if dialogwidth == "" {
		dialogwidth = "172pt"
	}
	return &DialogBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 0, Vexpanding: 0, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Title: Title, bOutSizeClose: bOutSizeClose, linkorbtn: linkorbtn, dialogwidth: dialogwidth}
}

//DialogBox
func (hb *DialogBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	return false
}
func (btn *DialogBox) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	var html string
	if btn.linkorbtn == 4 {
		if strings.Contains(btn.MaxPercentwidth, "%") {
			html = `<button` + id + ` onclick="document.location.href='#` + btn.id + `_openModal';d('` + btn.id + `_openModal').style.display='block';" style='width:` + btn.CalcWidth() + `;height:` + btn.CalcHeight() + `;` + btn.selfstyle + `'>` + btn.Title + `</button>`
		} else {
			html = `<button` + id + ` onclick="document.location.href='#` + btn.id + `_openModal';d('` + btn.id + `_openModal').style.display='block';" style='width:` + btn.MaxPercentwidth + `;height:` + btn.MaxPercentheight + `;` + btn.selfstyle + `'>` + btn.Title + `</button>`
		}
	} else {
		var textalign string
		if btn.linkorbtn == 2 {
			textalign = "text-align:center;"
		} else if btn.linkorbtn == 3 {
			textalign = "text-align:right;"
		} else {
			textalign = "text-align:left;"
		}
		//fmt.Println("btn.CalcWidth():", btn.CalcWidth(), btn.MaxPercentwidth)
		if strings.Contains(btn.MaxPercentwidth, "%") {
			html = `<div` + id + ` style='display:inline-block;width:auto;height:auto;` + textalign + `;` + btn.selfstyle + `'><a href='javascript:void(0)' onclick="document.location.href='#` + btn.id + `_openModal';d('` + btn.id + `_openModal').style.display='block';">` + btn.Title + `</a></div>`
		} else {
			html = `<div` + id + ` style='display:inline-block;width:auto;height:auto;` + textalign + `;` + btn.selfstyle + `'><a href='javascript:void(0)' onclick="document.location.href='#` + btn.id + `_openModal';d('` + btn.id + `_openModal').style.display='block';">` + btn.Title + `</a></div>`
		}
	}
	html += `<div id="` + btn.id + `_openModal" Class="modalDialog" style='display:none' >`
	if btn.bOutSizeClose {
		html += `<script type="text/javascript">
document.getElementById("` + btn.id + `_openModal").onclick=function (ev){
	var aaa=document.getElementById("dlg_` + btn.id + `")
	var style = null;
	if (window.getComputedStyle) {
		style = window.getComputedStyle(aaa, null);    // 非IE
	} else { 
		style = aaa.currentStyle;  // IE
	}
	ev=window.event||ev;  
	var scrollX = document.documentElement.scrollLeft || document.body.scrollLeft;  
	var scrollY = document.documentElement.scrollTop || document.body.scrollTop;  
	//var x = ev.pageX || ev.clientX + scrollX;  
	var y = ev.pageY || ev.clientY + scrollY;  
	//if(x<aaa.offsetLeft || y<aaa.offsetTop || x>aaa.offsetLeft+aaa.offsetWidth || y>aaa.offsetTop+aaa.offsetHeight)
	if(y<aaa.offsetTop || y>aaa.offsetTop+aaa.offsetHeight)
	{
		document.getElementById("` + btn.id + `_close").click();
		d('` + btn.id + `_openModal').style.display='none';
	}
}
	</script>`
	}
	html += `<div id="dlg_` + btn.id + `" style="width:` + btn.dialogwidth + `" >
            <a href="#` + btn.id + `_close" id="` + btn.id + `_close" title="Close" Class="modalDialogclose" onclick="d('` + btn.id + `_openModal').style.display='none';">X</a>
            <h3 style='margin-top:0pt;'>` + btn.Title + `</h3>
            __dialogitem__
        </div>
    </div>`
	var childstr string
	for _, val := range btn.child {
		str := val.ToHtml(countryid)
		childstr += str
	}
	return strings.Replace(html, "__dialogitem__", childstr, 1)
}

type ColorPickerDialog struct {
	BaseBox
	Title             string
	bOutSizeClose     bool
	linkorbtn         int
	dialogwidth, okjs string
}

//linkorbtn:1-left align link;2-center align link;3-right align link;4-button;getcolorvalue d(id+‘colorvalue').innerHTML;
func NewColorPickerDialog(id, Class string, color []string, Title, dialogwidth string, bOutSizeClose bool, linkorbtn int, okjs string) *ColorPickerDialog {
	if dialogwidth == "" {
		dialogwidth = "172pt"
	}
	return &ColorPickerDialog{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 0, Vexpanding: 0, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Title: Title, bOutSizeClose: bOutSizeClose, linkorbtn: linkorbtn, dialogwidth: dialogwidth, okjs: okjs}
}

//DialogBox
func (hb *ColorPickerDialog) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	return false
}
func (hb *ColorPickerDialog) ToHtml(countryid int) string {
	dlg := NewDialogBox(hb.id, hb.Class, hb.color, hb.Title, hb.dialogwidth, hb.bOutSizeClose, hb.linkorbtn)
	dlg.Add(NewOutHtml(hb.id+"colorpickerouthtml", "", nil, `
<form name="recherche" method="post" action="yourpage.html">
<input type=hidden name="rgb" value="123">
<table style="background-color:#f6f6f6;border:1px dotted #666;padding:5px;margin:0px auto;">
<tr>
<td style="border:1px outset #CCF;background-color:#ffe;width=172">
<div id=temoin style='float:right;width:40px;height:128px;'> </div>

<script language="Javascript" type="text/javascript">
<!--
var total=1657;var X=Y=j=RG=B=0;
var aR=new Array(total);var aG=new Array(total);var aB=new Array(total);
for (var i=0;i<256;i++){
aR[i+510]=aR[i+765]=aG[i+1020]=aG[i+5*255]=aB[i]=aB[i+255]=0;
aR[510-i]=aR[i+1020]=aG[i]=aG[1020-i]=aB[i+510]=aB[1530-i]=i;
aR[i]=aR[1530-i]=aG[i+255]=aG[i+510]=aB[i+765]=aB[i+1020]=255;
if(i<255){aR[i/2+1530]=127;aG[i/2+1530]=127;aB[i/2+1530]=127;}
}
function p(){var jla=document.getElementById('`+hb.id+`colorvalue');jla.innerHTML=artabus;jla.style.backgroundColor=artabus;document.forms['recherche'].rgb.value=artabus}
var hexbase=new Array("0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F");
var i=0;var jl=new Array();
for(x=0;x<16;x++)for(y=0;y<16;y++)jl[i++]=hexbase[x]+hexbase[y];
document.write('<'+'table border="0" cellspacing="0" cellpadding="0" onMouseover="t(event)" onClick="p()">');
var H=W=63;
for (Y=0;Y<=H;Y++){
	s='<'+'tr height=2>';j=Math.round(Y*(510/(H+1))-255)
	for (X=0;X<=W;X++){
		i=Math.round(X*(total/W))
		R=aR[i]-j;if(R<0)R=0;if(R>255||isNaN(R))R=255
		G=aG[i]-j;if(G<0)G=0;if(G>255||isNaN(G))G=255
		B=aB[i]-j;if(B<0)B=0;if(B>255||isNaN(B))B=255
		s=s+'<'+'td width=2 bgcolor=#'+jl[R]+jl[G]+jl[B]+'><'+'/td>'
	}
	document.write(s+'<'+'/tr>')
}
document.write('<'+'/table>');
var ns6=document.getElementById&&!document.all
var ie=document.all
var artabus=''
function t(e){
source=ie?event.srcElement:e.target
if(source.tagName=="TABLE")return
while(source.tagName!="TD" && source.tagName!="HTML")source=ns6?source.parentNode:source.parentElement
document.getElementById('temoin').style.backgroundColor=artabus=source.bgColor
}
// -->
</script>
<div id='`+hb.id+`colorvalue' style='height:24px;'  onClick="document.forms['recherche'].rgb.value='';this.style.backgroundColor=''"> </div><td></tr>
<tr><td colspan=2 align=center></td></tr>
</table>
</form>
	`))
	colorpickerokbtn := NewButton(hb.id+"colorpickerokbtn", "", nil, "Ok", hb.okjs)
	dlg.Add(colorpickerokbtn)
	return dlg.ToHtml(0)
}

type Label struct {
	BaseBox
	text      string
	textalign string
}

func NewLabel(id, Class string, color []string, text, textalign string) *Label {
	return &Label{BaseBox: BaseBox{id: id, Class: Class, color: color, Margin_left: "0pt", Margin_right: "0pt", Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, text: text, textalign: textalign}
}

//Label
func (btn *Label) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	return "<label" + id + " style='text-align:" + btn.textalign + ";width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";display:inline-block;white-space:nowrap;float:" + btn.Float + ";margin-left:" + btn.Margin_left + ";margin-right:" + btn.Margin_right + ";" + btn.selfstyle + "' " + btn.AttrToText() + ">" + btn.text + "</label>"
}

type Link struct {
	BaseBox
	Name, Onclickcode, Href, Target string
	textalign                       string
}

func NewLink(id, Class string, color []string, name, onclickcode, href, target, textalign string) *Link {
	return &Link{BaseBox: BaseBox{id: id, Class: Class, color: color, Margin_left: "0pt", Margin_right: "0pt", Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Onclickcode: onclickcode, Name: name, Href: href, Target: target, textalign: textalign}
}

//Link
func (btn *Link) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	hrefattr := ""
	if btn.Href != "" {
		hrefattr = "href='" + btn.Href + "'"
	}
	clkjs, scriptjs := getClickJsExpr(btn, btn.Onclickcode)
	return "<a" + id + " style='text-align:" + btn.textalign + ";width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";display:inline-block;float:" + btn.Float + ";margin-left:" + btn.Margin_left + ";margin-right:" + btn.Margin_right + ";" + btn.selfstyle + "' " + hrefattr + " target='" + btn.Target + "' onclick=\"" + clkjs + "\"" + btn.AttrToText() + ">" + btn.Name + "</a>" + scriptjs
}

type Input struct {
	BaseBox
	Type                    string //file checkbox text radio password hidden
	DefVal                  string
	PlaceHolder             string
	OnChangeJsCode, keyupjs string
}

//Input have parent.use parentNode to change something.radio:placeholder equal name;type:file:placeholder is the onchange js code;radio:placeholder is group of name;parent style use parentNode;
func NewInput(id, Class string, color []string, Type, DefVal, PlaceHolder, OnChangeJsCode, keyupjs string) *Input {
	return &Input{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Type: Type, DefVal: DefVal, PlaceHolder: PlaceHolder, OnChangeJsCode: OnChangeJsCode, keyupjs: keyupjs}
}

func (btn *Input) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	if btn.Type == "hidden" {
		return "<input" + id + " type='" + btn.Type + "' value='" + btn.DefVal + "'" + btn.AttrToText() + "/>"
	}
	pid := ""
	if btn.id != "" {
		pid = "id='" + btn.id + "p' "
	}
	if strings.Contains(btn.MaxPercentwidth, "%") {
		if btn.Type == "radio" {
			ck := ""
			if btn.DefVal[0] == '#' {
				ck = "checked='checked'"
				btn.DefVal = btn.DefVal[1:]
			}
			return "<div " + pid + "style='display:inline-block;width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";" + btn.selfstyle + "'><label><input" + id + " type='" + btn.Type + "' name='" + btn.PlaceHolder + "' " + ck + " style='display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' />" + btn.DefVal + "</label></div>"
		} else if btn.Type == "checkbox" {
			if btn.DefVal == "true" {
				return "<input" + id + " type='" + btn.Type + "' placeholder='" + btn.PlaceHolder + "' " + btn.DefVal + " style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' " + btn.AttrToText() + " checked/>"
			} else {
				return "<input" + id + " type='" + btn.Type + "' placeholder='" + btn.PlaceHolder + "' " + btn.DefVal + " style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' " + btn.AttrToText() + "/>"
			}
		} else if btn.Type == "file" {
			inljscode, scriptjscode := getClickJsExprWithEvent(btn, btn.OnChangeJsCode)
			keyupjslinejs, keyupjsscript := getClickJsExprWithEvent(btn, btn.keyupjs)
			return "<input" + id + " type='" + btn.Type + "' multiple='" + btn.DefVal + "' style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";" + btn.selfstyle + "' onchange=\"" + inljscode + "\" onkeyup=\"" + keyupjslinejs + "\"" + btn.AttrToText() + " />" + scriptjscode + keyupjsscript
		} else {
			inljscode, scriptjscode := getClickJsExprWithEvent(btn, btn.OnChangeJsCode)
			keyupjslinejs, keyupjsscript := getClickJsExprWithEvent(btn, btn.keyupjs)
			return "<input" + id + " type='" + btn.Type + "' placeholder='" + btn.PlaceHolder + "' value='" + btn.DefVal + "' style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + "display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";" + btn.selfstyle + "' onchange=\"" + inljscode + "\" onkeyup=\"" + keyupjslinejs + "\"" + btn.AttrToText() + " />" + scriptjscode + keyupjsscript
		}
	} else {
		if btn.Type == "radio" {
			ck := ""
			if btn.DefVal[0] == '#' {
				ck = "checked='checked'"
				btn.DefVal = btn.DefVal[1:]
			}
			return "<div " + pid + "style='display:inline-block;width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "'><label><input" + id + " type='" + btn.Type + "' name='" + btn.PlaceHolder + "' " + ck + " style='display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' />" + btn.DefVal + "</label></div>"
		} else if btn.Type == "checkbox" {
			if btn.DefVal == "true" {
				return "<div style='display:inline-block;width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "'><input" + id + " type='" + btn.Type + "' placeholder='" + btn.PlaceHolder + "' checked='" + btn.DefVal + "' style='width:100%;height:100%;display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' checked /></div>"
			} else {
				return "<div style='display:inline-block;width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "'><input" + id + " type='" + btn.Type + "' placeholder='" + btn.PlaceHolder + "' checked='" + btn.DefVal + "' style='width:100%;height:100%;display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' /></div>"
			}
		} else if btn.Type == "file" {
			inljscode, scriptjscode := getClickJsExprWithEvent(btn, btn.OnChangeJsCode)
			keyupjslinejs, keyupjsscript := getClickJsExprWithEvent(btn, btn.keyupjs)
			return "<div style='display:inline-block;width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "'><input" + id + " type='" + btn.Type + "' multiple='" + btn.DefVal + "' style='width:100%;height:100%;display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";'\"" + inljscode + "\" onkeyup=\"" + keyupjslinejs + "\" /></div>" + scriptjscode + keyupjsscript
		} else {
			inljscode, scriptjscode := getClickJsExprWithEvent(btn, btn.OnChangeJsCode)
			keyupjslinejs, keyupjsscript := getClickJsExprWithEvent(btn, btn.keyupjs)
			return "<div style='display:inline-block;width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "'><input" + id + " type='" + btn.Type + "' placeholder='" + btn.PlaceHolder + "' value='" + btn.DefVal + "' style='width:100%;height:100%;display:inline-block;border:1pt;padding:0pt;float:" + btn.Float + ";' onchange=\"" + inljscode + "\" onkeyup=\"" + keyupjslinejs + "\" /></div>" + scriptjscode + keyupjsscript
		}
	}
}

type AttrLabel struct {
	BaseBox
	attrname, value string
}

//Input have parent.use parentNode to change something.radio:placeholder equal name;type:file:placeholder is the onchange js code;radio:placeholder is group of name;parent style use parentNode;
func NewAttrLabel(id, Class string, color []string, attrname, value string) *AttrLabel {
	return &AttrLabel{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, attrname: attrname, value: value}
}

func (btn *AttrLabel) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	html := `<div><span><label>` + btn.attrname + `</label><label ` + id + ` >` + btn.value + `</label></span></div>`
	return html
}

type AttrInput struct {
	BaseBox
	attrname, value string
}

//Input have parent.use parentNode to change something.radio:placeholder equal name;type:file:placeholder is the onchange js code;radio:placeholder is group of name;parent style use parentNode;
func NewAttrInput(id, Class string, color []string, attrname, value string) *AttrInput {
	return &AttrInput{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, attrname: attrname, value: value}
}

func (btn *AttrInput) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	html := `<div><span><label>` + btn.attrname + `</label><input type='edit' ` + id + ` value='` + btn.value + `'/></span></div>`
	return html
}

type AttrNumber struct {
	BaseBox
	attrname, value string
}

//Input have parent.use parentNode to change something.radio:placeholder equal name;type:file:placeholder is the onchange js code;radio:placeholder is group of name;parent style use parentNode;
func NewAttrNumber(id, Class string, color []string, attrname, value string) *AttrNumber {
	return &AttrNumber{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, attrname: attrname, value: value}
}

func (btn *AttrNumber) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	html := `<div><span><label>` + btn.attrname + `</label><input type='number' ` + id + ` value='` + btn.value + `' /></span></div>`
	return html
}

type Button struct {
	BaseBox
	Text      string
	Onclickjs string
}

func NewButton(id, Class string, color []string, btnname, onclickjs string) *Button {
	return &Button{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Text: btnname, Onclickjs: strings.Trim(onclickjs, "\r\n\t ")}
}

//Button
func (btn *Button) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	rand.Seed(time.Now().UnixNano())
	var html string
	if strings.Contains(btn.MaxPercentwidth, "%") {
		html = "<button" + id + " style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";"
	} else {
		html = "<button" + id + " style='width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";"
	}

	clkjs, scriptjs := getClickJsExpr(btn, btn.Onclickjs)
	html += ";display:" + btn.Display + ";float:" + btn.Float + ";" + btn.selfstyle + "' onclick=\"" + clkjs + "\"" + btn.AttrToText() + ">" + btn.Text + "</button>" + scriptjs
	return html
}
func (btn *Button) OnClick(jscode string) {
	btn.Onclickjs = jscode
}

type EditArea struct {
	BaseBox
	EditHtml                        string
	uploadImageUrl, acceptextension string
}

//translate keys:FontSize,MakeHeader,TextAlignment,MakeLink,MakeList,InsertImage,Undo,Redo,Font,FontColor,BackgroundColor,OK,ImageFile,Keywords,Upload,Link
//upload image url param:keywords,imagefile;  accept eg:"image/jpeg, image/png"
func NewEditArea(id, Class string, color []string, EditHtml, uploadImageUrl, acceptextension string) *EditArea {
	return &EditArea{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 1, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, EditHtml: EditHtml, uploadImageUrl: uploadImageUrl, acceptextension: acceptextension}
}

//EditArea-
func (btn *EditArea) ToHtml(countryid int) string {
	colorddlg := NewColorPickerDialog(btn.id+"colorpickerdlg", "", nil, "Color Picker", "180pt", true, 4, `
	editarea['`+btn.id+`']['setTextColour'](d('`+btn.id+`colorpickerdlgcolorvalue').innerHTML);
	d('`+btn.id+`colorpickerdlg_openModal').style.display='none';
	`)
	colorpickerdlghtml := colorddlg.ToHtml(0)
	colorpickerdlghtml = colorpickerdlghtml[strings.Index(colorpickerdlghtml, "</button>")+len("</button>"):]
	//toolfunc.WriteFile("colordlg.txt", []byte(colorpickerdlghtml), 0666)
	//fmt.Println("colorpickerdlghtml", colorpickerdlghtml)

	var inputextension string
	if btn.acceptextension == "" {
		inputextension = ""
	} else {
		inputextension = btn.acceptextension
	}
	picuploadinput := NewInput(btn.id+"fileinput'  accept='"+inputextension+"'  style='display:none;", "", nil, "file", "", "", `
	if(d('`+btn.id+`fileinput').files.length>1){alert('`+T(T_PleaseSelectOneFile, countryid)+`');return;}
var form = new FormData();
form.append("editareafile", d('`+btn.id+`fileinput').files[0]);
postfile("`+btn.uploadImageUrl+`",form,function(ret){
	var retjn=JSON.parse(ret);
	if(retjn.error==0){
		editarea['`+btn.id+`']['insertImage'](retjn.imagelink,{"title":retjn.imagetitle});
	}
},function(progress){
	
});
	`, "")
	picuploadinputhtml := picuploadinput.ToHtml(0)

	squire_html := picuploadinputhtml + colorpickerdlghtml + `
  <style type="text/css" media="screen">
@font-face {
	font-family: FA;
	src: 
	url("/fa-solid-900.woff") format("woff");
	}
#` + btn.id + `alignment:before {content: "\f037";font-family: FA;}
#` + btn.id + `setTextColour:before {content: "\f031";font-family: FA;}
#` + btn.id + `setHighlightColour:before {content: "\f031";font-family: FA;}
#` + btn.id + `makeLink:before {content: "\f0c1";font-family: FA;}
#` + btn.id + `fontColor:before {content: "\f53f";font-family: FA;}
#` + btn.id + `makeHeader:before {content: "\f1dc";font-family: FA;}
#` + btn.id + `makeUnorderedList:before {content: "\f03a";font-family: FA;}
#` + btn.id + `insertFile:before {content: "\f15b";font-family: FA;}
#` + btn.id + `undo:before {content: "\f0e2";font-family: FA;}
#` + btn.id + `redo{transform:rotateY(180deg);font-family: FA;}
#` + btn.id + `redo:before {content: "\f0e2";font-family: FA;}

#` + btn.id + ` {
  -moz-box-sizing: border-box;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  min-height: 200px;
  border: 1px solid #888;
  padding: 1em;
  background: transparent;
  color: #2b2b2b;
  font: 13px/1.35 Helvetica, arial, sans-serif;
  cursor: text;
}
</style>


<header>
	<span id="` + btn.id + `makeHeader" class="edc" title="` + T(T_MAKEHEADER, countryid) + `"></span>
	<span id="` + btn.id + `alignment" class="edc" title="` + T(T_TEXTALIGNMENT, countryid) + `"></span>
  <span id="` + btn.id + `makeLink"  class="edc prompt" title="` + T(T_MAKELINK, countryid) + `"></span>
  <span id="` + btn.id + `fontColor"  class="edc prompt" title="` + T(T_FONTCOLOR, countryid) + `"></span>
  <span id="` + btn.id + `makeUnorderedList" class="edc" title="` + T(T_MAKELIST, countryid) + `"></span>
  <span id="` + btn.id + `insertFile"  class="edc prompt" title="` + T(T_INSERTFILE, countryid) + `"></span>
  <span id="` + btn.id + `undo" class="edc" title="` + T(T_UNDO, countryid) + `"></span>
  <div id="` + btn.id + `redo" class="edc" title="` + T(T_REDO, countryid) + `" style="display:inline-block;" ></div>
</header>
<div id="` + btn.id + `" style="overflow:auto;` + btn.selfstyle + `"></div>
<script type="text/javascript" charset="utf-8">
  var div = document.getElementById( '` + btn.id + `' );
  var editor = new Squire( div, {
      blockTag: 'p',
      blockAttributes: {'class': 'paragraph'},
      tagAttributes: {
          ul: {'class': 'UL'},
          ol: {'class': 'OL'},
          li: {'class': 'listItem'},
          a: {'target': '_blank'}
      }
  });
	editarea['` + btn.id + `']=editor;
  
  d("` + btn.id + `fontColor").onclick=function(){
document.location.href='#` + btn.id + `colorpickerdlg_openModal';d('` + btn.id + `colorpickerdlg_openModal').style.display='block';
};
var preselwidth=-1,preseltop='10pt',preseltop2=-1,curalignment;
d("` + btn.id + `alignment").onclick=function(){
	cpos=editarea['` + btn.id + `'].getCursorPosition();
	if(parseInt(cpos.width)!=preselwidth || cpos.top!=preseltop){
		preselwidth=parseInt(cpos.width);
		if(parseInt(preseltop)!=parseInt(cpos.top)){
			curalignment="left";
		}
		preseltop=cpos.top;
	}
	if(curalignment=="left"){
		curalignment="center";
	}else if(curalignment=="center"){
		curalignment="right";
	}else if(curalignment=="right"){
		curalignment="left";
	}
	editarea['` + btn.id + `']["setTextAlignment"](curalignment);
};
d("` + btn.id + `makeLink").onclick=function(){value = prompt('` + T(T_MAKELINK, countryid) + `');if(value=='')return;editarea['` + btn.id + `']['makeLink'](value);editarea['` + btn.id + `'].focus();};
var premakeheader;
d("` + btn.id + `makeHeader").onclick=function(){
if(preseltop=='10pt'){
editarea['` + btn.id + `']['setFontSize']('15pt');
preseltop='15pt';
}else if(preseltop=='15pt'){
editarea['` + btn.id + `']['setFontSize']('30pt');
preseltop='30pt';
}else if(preseltop=='30pt'){
editarea['` + btn.id + `']['setFontSize']('60pt');
preseltop='60pt';
}else if(preseltop=='60pt'){
editarea['` + btn.id + `']['setFontSize']('10pt');
preseltop='10pt';
}
editarea['` + btn.id + `'].focus();
};
var premakeUnorderedList;
d("` + btn.id + `makeUnorderedList").onclick=function(){
cpos=editarea['` + btn.id + `'].getCursorPosition();
if(parseInt(cpos.width)!=preselwidth || cpos.top!=preseltop2){
	preselwidth=parseInt(cpos.width);
	if(parseInt(cpos.top)!=parseInt(preseltop2)){
		premakeUnorderedList="removeList";
	}
	preseltop2=cpos.top;
}
if(premakeUnorderedList=="removeList"){
	premakeUnorderedList="makeUnorderedList";
}else{
	premakeUnorderedList="removeList";
}
editarea['` + btn.id + `'][premakeUnorderedList]('');
editarea['` + btn.id + `'].focus();
};
d("` + btn.id + `insertFile").onclick=function(){
d('` + btn.id + `fileinput').click();
};
d("` + btn.id + `undo").onclick=function(){editarea['` + btn.id + `']['undo']('');editarea['` + btn.id + `'].focus();};
d("` + btn.id + `redo").onclick=function(){editarea['` + btn.id + `']['redo']('');editarea['` + btn.id + `'].focus();};
`
	if btn.EditHtml != "" {
		squire_html += "editarea['" + btn.id + "'].setHTML(" + TextAsJavascriptString(btn.EditHtml) + ");"
	}
	squire_html += "</script>"

	return squire_html
}

func TextAsJavascriptString(text string) string {
	textls := regexp.MustCompile("[\r\n]+").Split(text, -1)
	re1 := regexp.MustCompile("\"")
	re2 := regexp.MustCompile("'")
	textstr := ""
	for i := 0; i < len(textls); i++ {
		if !re1.MatchString(text) {
			textstr += "\"" + text + "\"+\n"
		} else if !re2.MatchString(text) {
			textstr += "'" + text + "'+\n"
		} else {
			text = strings.Replace(text, "\"", "\\\"", -1)
			textstr += "\"" + text + "\"+\n"
		}
	}
	if strings.HasSuffix(textstr, "+\n") {
		textstr = textstr[:len(textstr)-2]
	}
	return textstr
}

type Combo struct {
	BaseBox
	Item           []string //file checkbox text radio password hidden
	Val            []string
	Name, changejs string
}

func NewCombo(id, Class string, color []string, Name, changejs string) *Combo {
	return &Combo{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 1, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Name: Name, changejs: changejs}
}

//Combo
func (hb *Combo) Add(name, value string) bool {
	//add to child list
	hb.Item = append(hb.Item, name)
	hb.Val = append(hb.Val, value)
	return false
}
func (btn *Combo) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	var itemstr string
	for i := 0; i < len(btn.Item); i++ {
		itemstr += "<option" + id + ` value ='` + btn.Val[i] + `'>` + btn.Item[i] + `</option>`
	}
	selecthtml := ""
	if strings.Contains(btn.MaxPercentwidth, "%") {
		selecthtml += "<select" + id + ` style='width:` + btn.CalcWidth() + `;height:` + btn.CalcHeight() + `;display:inline-block;float:` + btn.Float + `;` + btn.selfstyle + `' name='` + btn.Name + `'` + btn.AttrToText() + `>` + itemstr + `</select>`
	} else {
		selecthtml += "<select" + id + ` style='width:` + btn.MaxPercentwidth + `;height:` + btn.MaxPercentheight + `;display:inline-block;float:` + btn.Float + `;` + btn.selfstyle + `' name='` + btn.Name + `'` + btn.AttrToText() + `>` + itemstr + `</select>`
	}
	selecthtml += `<!--onloadbegin` + btn.changejs + `onloadend-->`
	return selecthtml
}

type VBox struct {
	BaseBox
	onclickjs string
}

func NewVBox(id, Class string, color []string, onclickjs string) *VBox {
	return &VBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, onclickjs: strings.Trim(onclickjs, "\r\n\t ")}
}

func (btn *VBox) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				hheight2 += maxhei
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *VBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	bc.SetDisplay("block")
	//set child width percent
	showcnt := 0
	for _, val := range hb.child {
		subhtml := val.ToHtml(Country_English)
		//fmt.Println("bh:", regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)), subhtml)
		if !regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)) {
			showcnt += 1
		}
	}
	onep := 100.0 / float32(showcnt)
	var ee float32
	var cnt int
	for _, val := range hb.child {
		subhtml := val.ToHtml(Country_English)
		if !regexp.MustCompile("^<[^>]*display:none;.*").Match([]byte(subhtml)) {
			if onep < val.GetVMaxPercent() {
				cnt += 1
			} else {
				val.SetVPercent(val.GetVMaxPercent())
				ee += onep - val.GetVMaxPercent()
			}
		}
	}
	onep += ee / float32(cnt)
	for _, val := range hb.child {
		if onep < val.GetVMaxPercent() {
			val.SetVPercent(onep)
		} else {

		}
	}

	return false
}

func (hb *VBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	html := "<div" + id + ` style="width:` + hb.CalcWidth() + `;height:auto;display:inline-block;vertical-align: middle;float:` + hb.Float + `;` + hb.selfstyle + `"`
	clkjs, scriptjs := getClickJsExpr(hb, hb.onclickjs)
	if clkjs != "" {
		html += " onclick=\"" + clkjs + "\""
	}
	html += hb.AttrToText() + ">__child__</div>"
	var childstr string
	for _, val := range hb.child {
		str := val.ToHtml(countryid)
		childstr += str
	}
	html = strings.Replace(html, "__child__", childstr, -1)
	if strings.Index(hb.onclickjs, "\n") != -1 {
		html += scriptjs
	}
	return html
}

type Pre struct {
	BaseBox
	Pretext string
}

func NewPre(id, Class string, color []string, pretext string) *Pre {
	return &Pre{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 3, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Pretext: pretext}
}

func (hb *Pre) Add(bc BaseControl) bool {
	return false
}
func (pr *Pre) ToHtml(countryid int) string {
	id := pr.ToIdAttr()
	html := "<pre" + id + " style='width:100%;height:auto;overflow:scroll;" + pr.selfstyle + "'" + pr.AttrToText() + ">" + pr.Pretext + "</pre>"
	return html
}

type P struct {
	BaseBox
	Ptext string
}

func NewP(id, Class string, color []string, ptext string) *P {
	return &P{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 3, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, Ptext: ptext}
}

func (hb *P) Add(bc BaseControl) bool {
	return false
}
func (pr *P) ToHtml(countryid int) string {
	id := pr.ToIdAttr()
	html := "<p" + id + " style='width:100%;height:auto;overflow:scroll;word-wrap:break-word;word-break:break-all;" + pr.selfstyle + "'" + pr.AttrToText() + ">" + pr.Ptext + "</p>"
	return html
}

//Bar
type Bar struct {
	BaseBox
	direct int
}

func NewBar(id, Class string, color []string, direct int) *Bar {
	//direct:up=1,right=2,bottom=3,left=4
	return &Bar{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 1, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, direct: direct}
}

func (hb *Bar) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	//set child width percent
	onep := 100.0 / float32(len(hb.child))
	var ee float32
	var cnt int
	for _, val := range hb.child {
		if onep < val.GetHMaxPercent() {
			cnt += 1
		} else {
			val.SetHPercent(val.GetHMaxPercent())
			ee += onep - val.GetHMaxPercent()
		}
	}
	onep += ee / float32(cnt)
	for _, val := range hb.child {
		if onep < val.GetHMaxPercent() {
			val.SetHPercent(onep)
		} else {

		}
	}
	return false
}

func (hb *Bar) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	var html string
	var childstr string
	if hb.direct == 1 || hb.direct == 3 {
		if hb.direct == 1 {
			html = "<div" + id + " class='topmostbar' style='vertical-align:middle;float:left;opacity:1;filter:alpha(opacity=100);-moz-opacity:1;opacity:1;border:0px;background:DeepSkyBlue;width:100%;position:fixed;top:0px;z-index:99999;clear:both;" + hb.selfstyle + "' >"
		} else if hb.direct == 3 {
			html = "<div" + id + " class='topmostbar' style='vertical-align:middle;float:left;opacity:1;filter:alpha(opacity=100);-moz-opacity:1;opacity:1;border:0px;background:DeepSkyBlue;width:100%;position:fixed;bottom:0px;z-index:99999;clear:both;" + hb.selfstyle + "' >"
		}
		for _, val := range hb.child {
			str := val.ToHtml(countryid)
			if len(str) > 5 && str[:5] == "<div " && !regexp.MustCompile("<div[^>]*display:inline-block;").Match([]byte(str)) {
				str = "<span " + str[4:len(str)-6] + "</span>"
			}
			childstr += str
		}
	} else {
		if hb.direct == 2 {
			html = "<div" + id + " class='topmostbar' style='align:center;float:left;opacity:1;filter:alpha(opacity=100);-moz-opacity:1;opacity:1;border:0px;background:DeepSkyBlue;height:100%;position:fixed;left:0px;z-index:99999;clear:both;" + hb.selfstyle + "' >"
		} else if hb.direct == 4 {
			html = "<div" + id + " class='topmostbar' style='align:center;:left;opacity:1;filter:alpha(opacity=100);-moz-opacity:1;opacity:1;border:0px;background:DeepSkyBlue;height:100%;position:fixed;right:0px;z-index:99999;clear:both;" + hb.selfstyle + "' >"
		}
		for _, val := range hb.child {
			str := val.ToHtml(countryid)
			childstr += "<div>" + str + "</div>"
		}
	}
	return html + childstr + "</div>"
}

//Checkbox
type CheckBox struct {
	BaseBox
	checkstate string
	text       string
	changejs   string
}

func NewCheckBox(id, Class string, color []string, checkstate, text, changejs string) *CheckBox {
	return &CheckBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, checkstate: checkstate, text: text, changejs: changejs}
}

func (btn *CheckBox) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	linejs, scriptjs := getClickJsExpr(btn, btn.changejs)
	html := "<a><input" + id + " type='checkbox' " + btn.checkstate + "='" + btn.checkstate + "' onchange='" + linejs + "' style='" + btn.selfstyle + "'" + btn.AttrToText() + " />" + btn.text + "</a>" + scriptjs
	return html
}

//radio group
type RadioGroup struct {
	BaseBox
	groupname       string
	defaultsel      string
	radiogroupnames []string
}

func NewRadioGroup(id, Class string, color []string, groupname, defaultsel string, radiogroupnames []string) *RadioGroup {
	return &RadioGroup{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, groupname: groupname, defaultsel: defaultsel, radiogroupnames: radiogroupnames}
}

func (btn *RadioGroup) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	rdstr := ""
	for _, name := range btn.radiogroupnames {
		if name != btn.defaultsel {
			rdstr += "<a><input" + id + " name='" + btn.groupname + "' type='radio' style='" + btn.selfstyle + "'" + btn.AttrToText() + "/>" + name + "</a>"
		} else {
			rdstr += "<a><input" + id + " name='" + btn.groupname + "' type='radio' checked='checked' style='" + btn.selfstyle + "'" + btn.AttrToText() + "/>" + name + "</a>"
		}
	}
	return rdstr
}

//image
type Image struct {
	BaseBox
	imagepath    string
	imagename    string
	onloadjscode string
}

func NewImage(id, Class string, color []string, imagepath, imagename, onloadjscode string) *Image {
	return &Image{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, imagepath: imagepath, imagename: imagename, onloadjscode: onloadjscode}
}

func (btn *Image) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	var html, onldseg string
	if btn.onloadjscode != "" {
		onldseg = "onload=\"" + btn.onloadjscode + "\""
	}
	if strings.Contains(btn.MaxPercentwidth, "%") {
		html = "<img" + id + " src='" + btn.imagepath + "' title='" + btn.imagename + "' style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";" + btn.selfstyle + "' " + onldseg + "" + btn.AttrToText() + "/>"
	} else {
		html = "<img" + id + "  src='" + btn.imagepath + "' title='" + btn.imagename + "' style='width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "' onload=\"" + btn.onloadjscode + "\" " + onldseg + "" + btn.AttrToText() + "/>"
	}
	return html
}

//textarea
type TextArea struct {
	BaseBox
	content string
}

func NewTextArea(id, Class string, color []string, content string) *TextArea {
	return &TextArea{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, content: content}
}

func (btn *TextArea) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	var html string
	if strings.Contains(btn.MaxPercentwidth, "%") {
		html = "<textarea" + id + " style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";" + btn.selfstyle + "' " + btn.AttrToText() + "/>" + btn.content + "</textarea>"
	} else {
		html = "<textarea" + id + " style='width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "' " + btn.AttrToText() + ">" + btn.content + "</textarea>"
	}
	return html
}

//imagetext
type ImageText struct {
	BaseBox
	imagepath string
	imagename string
	imagetext string
}

func NewImageText(id, Class string, color []string, imagepath, imagename, imagetext string) *ImageText {
	return &ImageText{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, imagepath: imagepath, imagename: imagename, imagetext: imagetext}
}

func (btn *ImageText) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	return "<div" + id + " style='" + btn.selfstyle + "'><div><img src='" + btn.imagepath + "' title='" + btn.imagename + "' /></div><div>" + btn.imagetext + "</div></div>"
}

type DropLoadBox struct {
	BaseBox
	AppendChildJsCode string
}

//AppendChildJsCode append child node at class id+lists with index 1
func NewDropLoadBox(id, Class string, color []string, AppendChildJsCode string) *DropLoadBox {
	return &DropLoadBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, AppendChildJsCode: AppendChildJsCode}
}

//DropLoadBox
func (hb *DropLoadBox) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	bc.SetFloat("left")
	return false
}

func (hb *DropLoadBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	html := "<div" + id + ` style="width:` + hb.Percentwidth + `;height:` + hb.Percentheight + `;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `"><div class="` + hb.id + `content"><div class="` + hb.id + `lists"></div><div class="` + hb.id + `lists"></div></div></div>
	<script type='text/javascript'>
	var result = '';
	var dropload = new DropLoad(d('` + hb.id + `'),{
        scrollArea : window,
		domUp:{
            domClass   : '` + hb.id + `dropload-up',
            domRefresh : '<div class="dropload-refresh">↓下拉刷新</div>',
            domUpdate  : '<div class="dropload-update">↑释放更新</div>',
            domLoad    : '<div class="dropload-load"><span class="loading"></span>加载中...</div>'
		},
		domDown:{
			domClass   : '` + hb.id + `dropload-down',
            domRefresh : '<div class="dropload-refresh">↑上拉加载更多</div>',
            domLoad    : '<div class="dropload-load"><span class="loading"></span>加载中...</div>',
            domNoData  : '<div class="dropload-noData">暂无数据</div>'
		},
        loadDownFn : function(me){
			` + hb.AppendChildJsCode + `
		}
	});
	droploadls['` + hb.id + `']=dropload;
	</script>
	`

	/*
			//sample append child code
			result +=   '<a class="item opacity" href="http://www.baidu.com/sdfsdf">'
		         +'<img src="http://www.g.cn/" alt="">'
		         +'<h3>493858sdjfklsdjafldsjf</h3>'
		         +'<span class="date">sdafdsafsafdaslkfj</span>'
		     +'</a>';
		    // 为了测试，延迟1秒加载
		    setTimeout(function(){
		        c('lists')[1].appendChild(h(result));
		        // 每次数据加载完，必须重置
		        me.resetload();
		    },1000);
	*/

	var childstr string
	for _, val := range hb.child {
		str := val.ToHtml(countryid)
		childstr += str
	}
	return strings.Replace(html, "__child__", childstr, -1)
}

//video
type Video struct {
	BaseBox
	videopath       string
	videoname       string
	oncanplayjscode string
}

func NewVideo(id, Class string, color []string, videopath, videoname, oncanplayjscode string) *Video {
	return &Video{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, videopath: videopath, videoname: videoname, oncanplayjscode: oncanplayjscode}
}

func (btn *Video) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	var html, onldseg string
	if btn.oncanplayjscode != "" {
		onldseg = "oncanplay=\"" + btn.oncanplayjscode + "\""
	}
	//oncanplay := "if(this.parentNode.offsetWidth-15>=this.videoWidth){this.style.width='auto';}else{this.style.width='100%';}"
	if strings.Contains(btn.MaxPercentwidth, "%") {
		html = "<video" + id + " src='" + btn.videopath + "' title='" + btn.videoname + "' oncanplay=\"" + btn.oncanplayjscode + "\" style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";" + btn.selfstyle + "' " + onldseg + ">Your browser does not support video.</video>"
	} else {
		html = "<video" + id + "  src='" + btn.videopath + "' title='" + btn.videoname + "' oncanplay=\"" + btn.oncanplayjscode + "\" style='width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "' " + onldseg + ">Your browser does not support video.</video>"
	}
	return html
}

//graphview
type GraphView struct {
	BaseBox
	videopath       string
	videoname       string
	oncanplayjscode string
}

func NewGraphView(id, Class string, color []string, videopath, videoname, oncanplayjscode string) *GraphView {
	return &GraphView{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, videopath: videopath, videoname: videoname, oncanplayjscode: oncanplayjscode}
}

func (btn *GraphView) ToHtml(countryid int) string {
	id := btn.ToIdAttr()
	var html, onldseg string
	if btn.oncanplayjscode != "" {
		onldseg = "oncanplay=\"" + btn.oncanplayjscode + "\""
	}
	//oncanplay := "if(this.parentNode.offsetWidth-15>=this.videoWidth){this.style.width='auto';}else{this.style.width='100%';}"
	if strings.Contains(btn.MaxPercentwidth, "%") {
		html = "<video" + id + " src='" + btn.videopath + "' title='" + btn.videoname + "' oncanplay=\"" + btn.oncanplayjscode + "\" style='width:" + btn.CalcWidth() + ";height:" + btn.CalcHeight() + ";" + btn.selfstyle + "' " + onldseg + "" + btn.AttrToText() + ">Your browser does not support video.</video>"
	} else {
		html = "<video" + id + "  src='" + btn.videopath + "' title='" + btn.videoname + "' oncanplay=\"" + btn.oncanplayjscode + "\" style='width:" + btn.MaxPercentwidth + ";height:" + btn.MaxPercentheight + ";" + btn.selfstyle + "' " + onldseg + "" + btn.AttrToText() + ">Your browser does not support video.</video>"
	}
	return html
}

//flowmain
type FlowMain struct {
	BaseBox
	videopath       string
	videoname       string
	oncanplayjscode string
}

func NewFlowMain(id, Class string, color []string) *FlowMain {
	return &FlowMain{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}}
}

func (btn *FlowMain) ToHtml(countryid int) string {
	idattr := btn.ToIdAttr()
	html := `<div ` + idattr + ` style="width:100%;height:auto;display:inline-block;float:none;overflow:auto;` + btn.selfstyle + `">
    __child__
</div>
<script type="text/javascript">
if(window.innerWidth/getDPI()>9){
document.getElementById('__flowmainid__').style.width="60%";
document.getElementById('__flowmainid__').style.marginLeft="20%";
document.getElementById('__flowmainid__').style.marginRight="20%";
ud=document.createElement("div");
ud.id="cornermax";ud.style.position="absolute";ud.style.zIndex="999";ud.style.float="left";ud.style.width="13pt";ud.style.height="13pt";ud.style.left="0pt";ud.style.top="0pt";ud.innerHTML="◤";
ud.onclick=function(){
   document.getElementById('cornermax') .style.display='none';
document.getElementById('__flowmainid__').style.width="100%";
document.getElementById('__flowmainid__').style.marginLeft="0%";
document.getElementById('__flowmainid__').style.marginRight="0%";
}
document.getElementById('__flowmainid__').parentElement.appendChild(ud);
}
</script>`
	html = strings.Replace(html, "__flowmainid__", btn.id, -1)
	var childstr string
	for _, val := range btn.child {
		str := val.ToHtml(countryid)
		childstr += str
	}
	return strings.Replace(html, "__child__", childstr, -1)
}

func (hb *FlowMain) Add(bc BaseControl) bool {
	//add to child list
	hb.child = append(hb.child, bc)
	bc.SetFloat("left")
	return false
}

type OutHtml struct {
	BaseBox
	outhtml  string
	bdivpack bool
}

func NewOutHtml(id, Class string, color []string, outhtml string) *OutHtml {
	return &OutHtml{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, outhtml: outhtml}
}

func (btn *OutHtml) SetDivPack(bdivpack bool) {
	btn.bdivpack = bdivpack
}

func (btn *OutHtml) ToHtml(countryid int) string {
	if btn.bdivpack {
		id := btn.ToIdAttr()
		return "<div " + id + ` style="` + btn.selfstyle + `">` + btn.outhtml + "</div>"
	} else {
		return btn.outhtml
	}
}

type Title struct {
	BaseBox
	title string
}

func NewTitle(id, Class string, color []string, title string) *Title {
	return &Title{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, title: title}
}

func (btn *Title) ToHtml(countryid int) string {
	return "<div style='text-align:center;font-weight=bold;" + btn.selfstyle + "'>" + btn.title + "</div>"
}

type Body struct {
	BaseBox
	body string
}

func NewBody(id, Class string, color []string, body string) *Body {
	return &Body{BaseBox: BaseBox{id: id, Class: Class, color: color, Maxheight: 20, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 1, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, body: body}
}

func (btn *Body) ToHtml(countryid int) string {
	return "<div style='" + btn.selfstyle + "'>" + btn.body + "</div>"
}

type Canvas2D struct {
	BaseBox
	onclickjs, mouseupjs, mousedownjs, mousemovejs string
}

func NewCanvas2D(id, Class string, color []string, onclickjs, mousedownjs, mouseupjs, mousemovejs string) *Canvas2D {
	return &Canvas2D{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, onclickjs: strings.Trim(onclickjs, "\r\n\t "), mouseupjs: mouseupjs, mousedownjs: mousedownjs, mousemovejs: mousemovejs}
}

func (btn *Canvas2D) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *Canvas2D) Add(bc BaseControl) bool {
	return false
}

func (hb *Canvas2D) ToHtml(countryid int) string {
	id := hb.ToIdAttr()
	html := "<canvas" + id + ` style="width:` + hb.CalcWidth() + `;height:` + hb.CalcHeight() + `;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `"`
	clkjs, scriptjs := getClickJsExpr(hb, hb.onclickjs)
	if clkjs != "" {
		html += " onclick=\"" + clkjs + "\""
	}
	html += `` + hb.AttrToText() + `>
	your browser didn't support canvas
	</canvas>`

	if hb.mouseupjs != "" || hb.mousedownjs != "" {
		html += "<!--onloadbegin"
		if hb.mousedownjs != "" {
			html += `d('` + hb.idname + `').addEventListener('mousedown', function(e) {
    let x = e.clientX  - d('` + hb.idname + `').getBoundingClientRect().left;
    let y = e.clientY  - d('` + hb.idname + `').getBoundingClientRect().top;
    ` + hb.mousedownjs + `
},true);`
		}
		if hb.mouseupjs != "" {
			html += `d('` + hb.idname + `').addEventListener('mouseup', function(e) {
    let x = e.clientX  - d('` + hb.idname + `').getBoundingClientRect().left;
    let y = e.clientY  - d('` + hb.idname + `').getBoundingClientRect().top;
    ` + hb.mouseupjs + `
},true);`
		}
		if hb.mousemovejs != "" {
			html += `d('` + hb.idname + `').addEventListener('mousemove', function(e) {
    let x = e.clientX  - d('` + hb.idname + `').getBoundingClientRect().left;
    let y = e.clientY  - d('` + hb.idname + `').getBoundingClientRect().top;
    ` + hb.mousemovejs + `
},true);`
		}
		html += "onloadend-->"
	}
	html += scriptjs

	return html
}

type ListBox struct {
	BaseBox
	onclickjs string
}

func NewListBox(id, Class string, color []string, onclickjs string) *ListBox {
	return &ListBox{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, onclickjs: strings.Trim(onclickjs, "\r\n\t ")}
}

func (btn *ListBox) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *ListBox) Add(bc BaseControl) bool {
	return false
}

func (hb *ListBox) ToHtml(countryid int) string {
	id := hb.ToIdAttr()

	html := "<select" + id + ` multiple='multiple'  size='12' style="width:` + hb.CalcWidth() + `;height:` + hb.CalcHeight() + `;display:inline-block;float:` + hb.Float + `;` + hb.selfstyle + `"`
	clkjs, scriptjs := getClickJsExpr(hb, hb.onclickjs)
	if clkjs != "" {
		html += " onclick=\"" + clkjs + "\""
	}
	html += `` + hb.AttrToText() + `>` + scriptjs + `
	</select>`
	return html
}

type Couplets struct {
	BaseBox
	leftname, leftlink, leftclickjs, rightname, rightlink, rightclickjs string
}

func NewCouplets(id, Class string, color []string, leftname, leftlink, leftclickjs, rightname, rightlink, rightclickjs string) *Couplets {
	return &Couplets{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, leftname: leftname, leftlink: leftlink, leftclickjs: leftclickjs, rightname: rightname, rightlink: rightlink, rightclickjs: rightclickjs}
}

func (btn *Couplets) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *Couplets) Add(bc BaseControl) bool {
	return false
}

func (hb *Couplets) ToHtml(countryid int) string {
	id := hb.ToIdAttr()

	html := "<div" + id + " style='" + hb.selfstyle + "'><div style='float:left'><a href='" + hb.leftlink + "'>" + hb.leftname + "</a></div><div style='float:right'><a href='" + hb.rightlink + "'>" + hb.rightname + "</a></div></div>"
	return html
}

type Carousel struct {
	BaseBox
	imagelistjs          []string
	firstimage, interval string
}

func NewCarousel(id, Class string, color []string) *Carousel {
	car := &Carousel{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}}
	car.interval = "3000"
	return car
}

func (btn *Carousel) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *Carousel) AddImageLink(link string) bool {
	hb.imagelistjs = append(hb.imagelistjs, link)
	return true
}

func (hb *Carousel) RemoveImageLink(link string) bool {
	ind := toolfunc.SliceSearch(hb.imagelistjs, link, 0)
	if ind != -1 {
		hb.imagelistjs = append(hb.imagelistjs[:ind], hb.imagelistjs[ind+1:]...)
	}
	return true
}

func (hb *Carousel) ToHtml(countryid int) string {
	idatr := hb.ToIdAttr()

	imagelistjsstr := ""
	numstr := ""
	firstimage := ""
	for i := 0; i < len(hb.imagelistjs); i++ {
		if firstimage == "" {
			firstimage = hb.imagelistjs[0]
		}
		imagelistjsstr += "'" + hb.imagelistjs[i] + "',"
		if numstr == "" {
			numstr = "<li class='selected'>1</li>"
		} else {
			numstr += "<li>" + toolfunc.IntToStr(i+1) + "</li>"
		}
	}
	if strings.HasSuffix(imagelistjsstr, ",") {
		imagelistjsstr = imagelistjsstr[:len(imagelistjsstr)-1]
	}
	if len(hb.imagelistjs) > 10 {
		numstr = ""
	}

	html := `<style type="text/css">
			body,
			ul,
			li {
				padding: 0;
				margin: 0;
			}
			
			.carouselBox {
				width: 512px;
				height: 384px;
				border: 3px solid black;
				margin: 0 auto;
				position: relative;
			}
			
			.carouselBox ul {
				position: absolute;
				left: 50%;
				transform: translateX(-50%);
				bottom: 10px;
				z-index: 2;
				overflow: hidden;
			}
			
			.carouselBox ul li {
				list-style: none;
				cursor: pointer;
				-moz-user-select: none;
				user-select: none;
				width: 18px;
				height: 18px;
				font-size: 14px;
				line-height: 18px;
				text-align: center;
				background: #ccc;
				float: left;
				margin: 2px;
			}
			
			.carouselBox ul li.selected {
				background: orange;
			}
			
			.carouselBox ul li.normal {
				background: #ccc;
			}
			
			.carouselBox button {
				position: absolute;
				top: 50%;
				transform: translateY(-50%);
				font-size: 40px;
				font-weight: 200;
				opacity: 0;
			}
			
			.carouselBox button.show {
				opacity: .7;
			}
			
			.carouselBox .btnL {
				left: 10px;
			}
			
			.carouselBox .btnR {
				right: 10px;
			}
		</style>
		
		<div ` + idatr + ` class="carouselBox" >
			<img class="carouselImg" src="` + firstimage + `" alt="" style="width:100%;height:100%;" />
			<ul>
				` + numstr + `
			</ul>
			<button class="btnL">&lt;</button>
			<button class="btnR">&gt;</button>
		</div>
		<script type="text/javascript">
			//封装轮播
			function carousel(carouselBox) {

				var carouselBox = carouselBox;
				//获取轮播图片节点
				var carouselImg = carouselBox.getElementsByClassName("carouselImg")[0];
				//获取到所以轮播按钮节点
				var lis = carouselBox.getElementsByTagName("li");
				var len = lis.length;
				//左点击按钮
				var btnL = carouselBox.getElementsByClassName("btnL")[0];
				//右点击按钮
				var btnR = carouselBox.getElementsByClassName("btnR")[0];
				//定义默认图片
				var autoImgList=[__AutoImageList__];
				var autoImg = 1;
				
				//定时器
				var time = null;
				//当前li下标值
				var currIndex = null;
				//调用自动轮播
				autoCarousel();
				//调用点击轮播
				clickCarousel();

				//自动轮播
				function autoCarousel() {
					time = setInterval(autoChange,` + hb.interval + `);
				}
				//自动切换函数
				function autoChange() {
					//如果是最后一张变成第一张
					if(autoImg+1 == autoImgList.length) {
						autoImg = 0;
					} else {
						autoImg++;
					}
					carouselImg.src = autoImgList[autoImg];

					//对应得按钮背景色改变
					bgChange(autoImg);
				}

				//背景色改变函数
				function bgChange(t) {

					for(var i = 0; i < len; i++) {
						if(i === t) {
							//如果是和当前图片对应得下标的li则改变背景 否则变为正常的
							lis[i].className = "selected";
						} else {
							lis[i].className = "normal";
						}
					}
				}

				//点击轮播
				function clickCarousel() {
					for(var i = 0; i < len; i++) {
						lis[i].onclick = function() {
							currIndex = getIndex(this);
							bgChange(currIndex);
							carouselImg.src = autoImgList[currIndex];
							autoImg = currIndex;
						}
					}
				}

				//获取到点击li的下标
				function getIndex(t) {
					//定义个标签
					var index = -1;
					for(var i = 0; i < len; i++) {
						//找到当前点击对象 并记录下标值
						if(lis[i] === t) {
							index = i;
							break;
						}
					}
					return index;
				}
				//鼠标移入的时候清楚定时器
				carouselBox.onmouseenter = function() {
					clearTimeout(time);
					btnL.className = "show btnL";
					btnR.className = "show btnR";
				}
				//鼠标移出继续轮播
				carouselBox.onmouseleave = function() {
					autoCarousel();
					btnL.className = "btnL";
					btnR.className = "btnR";
				}
				//左点击按钮事件
				btnL.onclick = function() {
					currIndex = autoImg;
					if(autoImg == 0) {
						autoImg = autoImgList.length-1;
					} else {
						autoImg--;
					}
					carouselImg.src = autoImgList[autoImg];

					//对应得按钮背景色改变
					bgChange(autoImg);
				}
				//右按钮点击事件
				btnR.onclick = function() {
					currIndex = autoImg;
					autoChange();
				}
			}
			//调用轮播
			carousel(document.getElementById("` + hb.id + `"));
			
		</script>
	`

	html = strings.ReplaceAll(html, "__AutoImageList__", imagelistjsstr)
	return html
}

type H3 struct {
	BaseBox
	title string
}

func NewH3(id, Class string, color []string, title string) *H3 {
	car := &H3{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, title: title}
	return car
}

func (btn *H3) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *H3) ToHtml(countryid int) string {
	idatr := hb.ToIdAttr()
	return "<h3 " + idatr + ">" + hb.title + "</h3>"
}

type H1 struct {
	BaseBox
	title string
}

func NewH1(id, Class string, color []string, title string) *H1 {
	car := &H1{BaseBox: BaseBox{id: id, Class: Class, color: color, Percentwidth: "100%", Percentheight: "100%", Hexpanding: 5, Vexpanding: 5, MinPercentwidth: "0%", MinPercentheight: "0%", MaxPercentwidth: "100%", MaxPercentheight: "100%", Float: "none", Display: "inline-block"}, title: title}
	return car
}

func (btn *H1) CalcHeight() string {
	if btn.Vexpanding == 1 {
		return strconv.FormatInt(int64(btn.Maxheight), 10) + "pt"
	} else if btn.Vexpanding == 3 {
		return "auto"
	} else {
		hheight := ""
		ballfixhei := true
		for _, val := range btn.child {
			if val.GetVExpanding() != 1 {
				if val.GetVExpanding() == 3 {
					return "auto"
				} else {
					ballfixhei = false
				}
			}
		}
		if ballfixhei == true {
			hheight2 := 0
			for _, val := range btn.child {
				maxhei := val.GetMaxHeight()
				//fmt.Println("hbox maxhei", maxhei)
				if maxhei > hheight2 {
					hheight2 = maxhei
				}
			}
			if hheight2 > 0 {
				hheight = strconv.FormatInt(int64(hheight2), 10) + "pt"
			}
		}
		if hheight != "" {
			return hheight
		} else {
			return strconv.FormatFloat(math.Min(percentv(btn.Percentheight), percentv(btn.MaxPercentheight)), 'f', 4, 32) + "%"
		}
	}
}

func (hb *H1) ToHtml(countryid int) string {
	idatr := hb.ToIdAttr()
	return "<h1 " + idatr + hb.AttrToText() + ">" + hb.title + "</h1>"
}
