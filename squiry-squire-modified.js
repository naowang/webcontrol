//below insertImage add
proto.insertAudio = function ( src, attributes ) {
    var img = this.createElement( 'audio', mergeObjects({
        src: src,controls:"controls"
    }, attributes, true ));
    this.insertElement( img );
    return img;
};

proto.insertVideo = function ( src, attributes ) {
    var img = this.createElement( 'video', mergeObjects({
        src: src,controls:"controls"
    }, attributes, true ));
    this.insertElement( img );
    return img;
};

proto.insertFile = function ( src, attributes ) {
	var fn=src.substr(src.lastIndexOf("/")+1);
    var img = this.createElement( 'a', mergeObjects({
        href: src, target: "_blank", download:fn
    }, attributes, true ));
	img.innerText=fn;
    this.insertElement( img );
    return img;
};