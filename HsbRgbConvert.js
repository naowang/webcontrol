function Rgb2Hsb(rgb) {
	hsb = Array();
	var maxIndex, minIndex;
	var tmp;
	rearranged = rgb.slice(0);
	for(i= 0;i<2;i++){
		for(j=0;j<2-i;j++){
			if(rearranged[j]>rearranged[j+1]){
				tmp=rearranged[j+1];
				rearranged[j+1]=rearranged[j];
				rearranged[j]=tmp;
			}
		}
	}
	for(i=0;i< 3; i++){
		if(rearranged[0]==rgb[i]){
			minIndex = i;
		}
		if(rearranged[2] == rgb[i]){
			maxIndex = i;
		}
	}
	hsb[2] = rearranged[2] / 255.0;
	if(rearranged[2] > 0){
		hsb[1] = 1 - rearranged[0]/rearranged[2];
	} else {
		hsb[1] = 0;
	}
	if(hsb[1] > 0){
		if((maxIndex-minIndex+3)%3 == 1){
			hsb[0] = parseFloat(maxIndex)*120 + 60*(rearranged[1]/hsb[1]/rearranged[2]+(1-1/hsb[1]))*1;
		}else{
			hsb[0] = parseFloat(maxIndex)*120 + 60*(rearranged[1]/hsb[1]/rearranged[2]+(1-1/hsb[1]))*-1;
		}
	} else {
		hsb[0] = 360;
	}
	console.log(hsb[0]);
	hsb[0] = parseFloat(parseInt(hsb[0]+360) % 360);
	return hsb;
}

function Hsb2Rgb(hsb) {
	rgb=Array();
	offset=240
	i=0
	for(;i<3;){
		x=parseFloat(Math.abs(parseFloat(parseInt(hsb[0]+parseFloat(offset))%360-240)));
		if(x <= 60){
			rgb[i]=255;
		}else if(60<x&&x<120){
			rgb[i]=((1-(x-60)/60)*255);
		}else{
			rgb[i]=0;
		}
		i++;
		offset-= 120;
	}
	for(i=0;i<3;i++){
		rgb[i]+=(255-rgb[i])*(1-hsb[1]);
	}
	for(i=0;i<3;i++){
		rgb[i]*=hsb[2];
	}
	return rgb;
}



