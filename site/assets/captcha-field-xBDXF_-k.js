import{i as e,o as t,r as n,t as r}from"./jsx-runtime-C1p6_exl.js";import{n as i}from"./client-Bb93rIye.js";import{n as a}from"./provider-BuZ6yQ3Y.js";import"./label-BwzFW-Ds.js";import"./button-BImRTnkK.js";import"./input-rWbN_PRx.js";var o=e(((e,t)=>{(function(){var e={}.hasOwnProperty;function n(){for(var e=``,t=0;t<arguments.length;t++){var n=arguments[t];n&&(e=i(e,r(n)))}return e}function r(t){if(typeof t==`string`||typeof t==`number`)return t;if(typeof t!=`object`)return``;if(Array.isArray(t))return n.apply(null,t);if(t.toString!==Object.prototype.toString&&!t.toString.toString().includes(`[native code]`))return t.toString();var r=``;for(var a in t)e.call(t,a)&&t[a]&&(r=i(r,a));return r}function i(e,t){return t?e?e+` `+t:e+t:e}t!==void 0&&t.exports?(n.default=n,t.exports=n):typeof define==`function`&&typeof define.amd==`object`&&define.amd?define(`classnames`,[],function(){return n}):window.classNames=n})()})),s=t(n()),c=t(o());function l(e,t){t===void 0&&(t={});var n=t.insertAt;if(!(!e||typeof document>`u`)){var r=document.head||document.getElementsByTagName(`head`)[0],i=document.createElement(`style`);i.type=`text/css`,n===`top`&&r.firstChild?r.insertBefore(i,r.firstChild):r.appendChild(i),i.styleSheet?i.styleSheet.cssText=e:i.appendChild(document.createTextNode(e))}}var u=`/**
 * @Author Awen
 * @Date 2024/06/01
 * @Email wengaolng@gmail.com
 **/
.index-module_iconBlock__Y1IUb {
  flex: 1;
}
.index-module_dots__2OJFw {
  position: absolute;
  top: 0;
  right: 0;
  left: 0;
  bottom: 0;
}
.index-module_dots__2OJFw .dot {
  position: absolute;
  z-index: 2;
  width: 22px;
  height: 22px;
  color: var(--go-captcha-theme-dot-color);
  background: var(--go-captcha-theme-dot-bg-color);
  border: 3px solid #f7f9fb;
  border-color: var(--go-captcha-theme-dot-border-color);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 22px;
  cursor: default;
}
`,d={iconBlock:`index-module_iconBlock__Y1IUb`,dots:`index-module_dots__2OJFw`};l(u);var f=`/**
 * @Author Awen
 * @Date 2024/06/01
 * @Email wengaolng@gmail.com
 **/
:root {
  --go-captcha-theme-text-color: #333333;
  --go-captcha-theme-bg-color: #ffffff;
  --go-captcha-theme-btn-color: #ffffff;
  --go-captcha-theme-btn-disabled-color: #749ff9;
  --go-captcha-theme-btn-bg-color: #4e87ff;
  --go-captcha-theme-btn-border-color: #4e87ff;
  --go-captcha-theme-active-color: #3e7cff;
  --go-captcha-theme-border-color: rgba(206, 223, 254, 0.5);
  --go-captcha-theme-icon-color: #3C3C3C;
  --go-captcha-theme-drag-bar-color: #e0e0e0;
  --go-captcha-theme-drag-bg-color: #3e7cff;
  --go-captcha-theme-drag-icon-color: #ffffff;
  --go-captcha-theme-round-color: #e0e0e0;
  --go-captcha-theme-loading-icon-color: #3e7cff;
  --go-captcha-theme-body-bg-color: #34383e;
  --go-captcha-theme-dot-color: #cedffe;
  --go-captcha-theme-dot-bg-color: #3e7cff;
  --go-captcha-theme-dot-border-color: #f7f9fb;
  --go-captcha-theme-default-color: #3e7cff;
  --go-captcha-theme-default-bg-color: #ecf5ff;
  --go-captcha-theme-default-border-color: #3e7cff;
  --go-captcha-theme-default-hover-color: #e0efff;
  --go-captcha-theme-error-color: #ed4630;
  --go-captcha-theme-error-bg-color: #fef0f0;
  --go-captcha-theme-error-border-color: #ff5a34;
  --go-captcha-theme-warn-color: #ffa000;
  --go-captcha-theme-warn-bg-color: #fdf6ec;
  --go-captcha-theme-warn-border-color: #ffbe09;
  --go-captcha-theme-success-color: #5eaa2f;
  --go-captcha-theme-success-bg-color: #f0f9eb;
  --go-captcha-theme-success-border-color: #8bc640;
}
.gocaptcha-module_wrapper__Kpdey {
  padding: 12px 16px;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
  box-sizing: border-box;
}
.gocaptcha-module_theme__h-Ytl {
  border: 1px solid rgba(206, 223, 254, 0.5);
  border-color: var(--go-captcha-theme-border-color);
  border-radius: 8px;
  box-shadow: 0 0 20px rgba(100, 100, 100, 0.1);
  -webkit-box-shadow: 0 0 20px rgba(100, 100, 100, 0.1);
  -moz-box-shadow: 0 0 20px rgba(100, 100, 100, 0.1);
  background-color: var(--go-captcha-theme-bg-color);
}
.gocaptcha-module_header__LjDUC {
  height: 36px;
  width: 100%;
  font-size: 15px;
  color: var(--go-captcha-theme-text-color);
  display: flex;
  align-items: center;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
.gocaptcha-module_header__LjDUC span {
  flex: 1;
  padding-right: 5px;
}
.gocaptcha-module_header__LjDUC em {
  padding: 0 3px;
  font-weight: bold;
  color: var(--go-captcha-theme-active-color);
  font-style: normal;
}
.gocaptcha-module_body__KJKNu {
  position: relative;
  width: 100%;
  margin-top: 10px;
  display: flex;
  background: var(--go-captcha-theme-body-bg-color);
  border-radius: 5px;
  -webkit-border-radius: 5px;
  -moz-border-radius: 5px;
  overflow: hidden;
}
.gocaptcha-module_bodyInner__jahqH {
  position: relative;
  background: var(--go-captcha-theme-body-bg-color);
}
.gocaptcha-module_picture__LRwbY {
  position: relative;
  z-index: 2;
  width: 100%;
}
.gocaptcha-module_hide__TUOZE {
  visibility: hidden;
}
.gocaptcha-module_loading__Y-PYK {
  position: absolute;
  z-index: 1;
  top: 50%;
  left: 50%;
  width: 68px;
  height: 68px;
  margin-left: -34px;
  margin-top: -34px;
  line-height: 68px;
  text-align: center;
  display: flex;
  align-content: center;
  justify-content: center;
}
.gocaptcha-module_loading__Y-PYK svg,
.gocaptcha-module_loading__Y-PYK circle {
  color: var(--go-captcha-theme-loading-icon-color);
  fill: var(--go-captcha-theme-loading-icon-color);
}
.gocaptcha-module_footer__Ywdpy {
  width: 100%;
  height: 50px;
  color: #34383e;
  display: flex;
  align-items: center;
  padding-top: 10px;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
.gocaptcha-module_iconBlock__mVB8B {
  display: flex;
  align-items: center;
}
.gocaptcha-module_iconBlock__mVB8B svg {
  color: var(--go-captcha-theme-icon-color);
  fill: var(--go-captcha-theme-icon-color);
  margin: 0 5px;
  cursor: pointer;
}
.gocaptcha-module_buttonBlock__EZ4vg {
  width: 120px;
  height: 40px;
}
.gocaptcha-module_buttonBlock__EZ4vg button {
  width: 100%;
  height: 40px;
  text-align: center;
  padding: 9px 15px;
  font-size: 15px;
  border-radius: 5px;
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
  cursor: pointer;
  color: var(--go-captcha-theme-btn-color);
  background-color: var(--go-captcha-theme-btn-bg-color);
  border: 1px solid transparent;
  border-color: var(--go-captcha-theme-btn-bg-color);
  -webkit-appearance: none;
  box-sizing: border-box;
  outline: none;
  margin: 0;
  transition: 0.1s;
  font-weight: 500;
  -moz-user-select: none;
  -webkit-user-select: none;
}
.gocaptcha-module_buttonBlock__EZ4vg button.disabled {
  pointer-events: none;
  background-color: var(--go-captcha-theme-btn-disabled-color);
  border-color: var(--go-captcha-theme-btn-disabled-color);
}
.gocaptcha-module_dragSlideBar__noauW {
  width: 100%;
  height: 100%;
  position: relative;
  touch-action: none;
}
.gocaptcha-module_dragLine__3B9KR {
  position: absolute;
  height: 14px;
  background-color: var(--go-captcha-theme-drag-bar-color);
  left: 0;
  right: 0;
  top: 50%;
  margin-top: -7px;
  border-radius: 7px;
}
.gocaptcha-module_dragBlock__bFlwx {
  position: absolute;
  left: 0;
  top: 50%;
  margin-top: -20px;
  width: 82px;
  height: 40px;
  z-index: 2;
  background-color: var(--go-captcha-theme-drag-bg-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
  border-radius: 24px;
  box-shadow: 0 0 20px rgba(100, 100, 100, 0.35);
  -webkit-box-shadow: 0 0 20px rgba(100, 100, 100, 0.35);
  -moz-box-shadow: 0 0 20px rgba(100, 100, 100, 0.35);
}
.gocaptcha-module_dragBlock__bFlwx svg {
  color: var(--go-captcha-theme-drag-icon-color);
  fill: var(--go-captcha-theme-drag-icon-color);
}
.gocaptcha-module_disabled__4kN6w {
  pointer-events: none;
  background-color: var(--go-captcha-theme-btn-disabled-color);
  border-color: var(--go-captcha-theme-btn-disabled-color);
}
.gocaptcha-module_dragBlockInline__PpF3f {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
`,p={wrapper:`gocaptcha-module_wrapper__Kpdey`,theme:`gocaptcha-module_theme__h-Ytl`,header:`gocaptcha-module_header__LjDUC`,body:`gocaptcha-module_body__KJKNu`,bodyInner:`gocaptcha-module_bodyInner__jahqH`,picture:`gocaptcha-module_picture__LRwbY`,hide:`gocaptcha-module_hide__TUOZE`,loading:`gocaptcha-module_loading__Y-PYK`,footer:`gocaptcha-module_footer__Ywdpy`,iconBlock:`gocaptcha-module_iconBlock__mVB8B`,buttonBlock:`gocaptcha-module_buttonBlock__EZ4vg`,dragSlideBar:`gocaptcha-module_dragSlideBar__noauW`,dragLine:`gocaptcha-module_dragLine__3B9KR`,dragBlock:`gocaptcha-module_dragBlock__bFlwx`,disabled:`gocaptcha-module_disabled__4kN6w`,dragBlockInline:`gocaptcha-module_dragBlockInline__PpF3f`};l(f);var m=()=>({width:300,height:220,thumbWidth:150,thumbHeight:40,verticalPadding:16,horizontalPadding:12,showTheme:!0,title:`请在下图依次点击`,buttonText:`确认`,iconSize:22,dotSize:24}),h=e=>(0,s.createElement)(`svg`,Object.assign({xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 200 200`,width:26,height:26},e),(0,s.createElement)(`path`,{d:`M100.1,189.9C100.1,189.9,100,189.9,100.1,189.9c-49.7,0-90-40.4-90-89.9c0-49.6,40.4-89.9,89.9-89.9
		c49.6,0,89.9,40.4,89.9,89.9c0,18.2-5.4,35.7-15.6,50.7c-1.5,2.1-3.6,3.4-6.1,3.9c-2.5,0.4-5-0.1-7-1.6c-4.2-3-5.3-8.6-2.4-12.9
		c8.1-11.9,12.4-25.7,12.4-40.1c0-39.2-31.9-71.1-71.1-71.1c-39.2,0-71.1,31.9-71.1,71.1c0,39.2,31.9,71.1,71.1,71.1
		c7.7,0,15.3-1.2,22.6-3.6c2.4-0.8,4.9-0.6,7.2,0.5c2.2,1.1,3.9,3.1,4.7,5.5c1.6,4.9-1,10.2-5.9,11.9
		C119.3,188.4,109.8,189.9,100.1,189.9z M73,136.4C73,136.4,73,136.4,73,136.4c-2.5,0-4.9-1-6.7-2.8c-3.7-3.7-3.7-9.6,0-13.3
		L86.7,100L66.4,79.7c-3.7-3.7-3.7-9.6,0-13.3c3.7-3.7,9.6-3.7,13.3,0L100,86.7l20.3-20.3c1.8-1.8,4.1-2.8,6.7-2.8c0,0,0,0,0,0
		c2.5,0,4.9,1,6.7,2.8c1.8,1.8,2.8,4.1,2.8,6.7c0,2.5-1,4.9-2.8,6.7L113.3,100l20.3,20.3c3.7,3.7,3.7,9.6,0,13.3
		c-3.7,3.7-9.6,3.7-13.3,0L100,113.3l-20.3,20.3C77.9,135.4,75.5,136.4,73,136.4z`})),g=e=>(0,s.createElement)(`svg`,Object.assign({width:26,height:26},e,{viewBox:`0 0 200 200`,xmlns:`http://www.w3.org/2000/svg`}),(0,s.createElement)(`path`,{d:`M135,149.9c-10.7,7.6-23.2,11.4-36,11.2c-1.7,0-3.4-0.1-5-0.3c-0.7-0.1-1.4-0.2-2-0.3c-1.3-0.2-2.6-0.4-3.9-0.6
	c-0.8-0.2-1.6-0.4-2.3-0.5c-1.2-0.3-2.5-0.6-3.7-1c-0.6-0.2-1.2-0.4-1.7-0.6c-1.4-0.5-2.8-1-4.2-1.5c-0.3-0.1-0.6-0.3-0.9-0.4
	c-1.6-0.7-3.2-1.4-4.7-2.3c-0.1,0-0.1-0.1-0.2-0.1c-5.1-2.9-9.8-6.4-14-10.6c-0.1-0.1-0.1-0.1-0.2-0.2c-1.3-1.3-2.5-2.7-3.7-4.1
	c-0.2-0.3-0.5-0.6-0.7-0.9c-8.4-10.6-13.5-24.1-13.5-38.8h14.3c0.4,0,0.7-0.2,0.9-0.5c0.2-0.3,0.2-0.8,0-1.1L29.5,60.9
	c-0.2-0.3-0.5-0.5-0.9-0.5c-0.4,0-0.7,0.2-0.9,0.5L3.8,97.3c-0.2,0.3-0.2,0.7,0,1.1c0.2,0.3,0.5,0.5,0.9,0.5h14.3
	c0,17.2,5.3,33.2,14.3,46.4c0.1,0.2,0.2,0.4,0.3,0.6c0.9,1.4,2,2.6,3,3.9c0.4,0.5,0.7,1,1.1,1.5c1.5,1.8,3,3.5,4.6,5.2
	c0.2,0.2,0.3,0.3,0.5,0.5c5.4,5.5,11.5,10.1,18.2,13.8c0.2,0.1,0.3,0.2,0.5,0.3c1.9,1,3.9,2,5.9,2.9c0.5,0.2,1,0.5,1.5,0.7
	c1.7,0.7,3.5,1.3,5.2,1.9c0.8,0.3,1.7,0.6,2.5,0.8c1.5,0.5,3.1,0.8,4.7,1.2c1.1,0.2,2.1,0.5,3.2,0.7c0.4,0.1,0.9,0.2,1.3,0.3
	c1.5,0.3,3,0.4,4.5,0.6c0.5,0.1,1.1,0.2,1.6,0.2c2.7,0.3,5.4,0.4,8.1,0.4c16.4,0,32.5-5.1,46.2-14.8c4.4-3.1,5.5-9.2,2.4-13.7
	C145.5,147.8,139.4,146.7,135,149.9 M180.6,98.9c0-17.2-5.3-33.1-14.2-46.3c-0.1-0.2-0.2-0.5-0.4-0.7c-1.1-1.6-2.3-3.1-3.5-4.6
	c-0.1-0.2-0.3-0.4-0.4-0.6c-8.2-10.1-18.5-17.9-30.2-23c-0.3-0.1-0.6-0.3-1-0.4c-1.9-0.8-3.8-1.5-5.7-2.1c-0.7-0.2-1.4-0.5-2.1-0.7
	c-1.7-0.5-3.4-0.9-5.1-1.3c-0.9-0.2-1.9-0.5-2.8-0.7c-0.5-0.1-0.9-0.2-1.4-0.3c-1.3-0.2-2.6-0.3-3.8-0.5c-0.9-0.1-1.8-0.3-2.6-0.3
	c-2.1-0.2-4.3-0.3-6.4-0.3c-0.4,0-0.8-0.1-1.2-0.1c-0.1,0-0.1,0-0.2,0c-16.4,0-32.4,5-46.2,14.8C49,35,48,41.1,51,45.6
	c3.1,4.4,9.1,5.5,13.5,2.4c10.6-7.5,23-11.3,35.7-11.2c1.8,0,3.6,0.1,5.4,0.3c0.6,0.1,1.1,0.1,1.6,0.2c1.5,0.2,2.9,0.4,4.3,0.7
	c0.6,0.1,1.3,0.3,1.9,0.4c1.4,0.3,2.8,0.7,4.2,1.1c0.4,0.1,0.9,0.3,1.3,0.4c1.6,0.5,3.1,1.1,4.6,1.7c0.2,0.1,0.3,0.1,0.5,0.2
	c9,3.9,17,10,23.2,17.6c0,0,0.1,0.1,0.1,0.2c8.7,10.7,14,24.5,14,39.4H147c-0.4,0-0.7,0.2-0.9,0.5c-0.2,0.3-0.2,0.8,0,1.1l24,36.4
	c0.2,0.3,0.5,0.5,0.9,0.5c0.4,0,0.7-0.2,0.9-0.5l23.9-36.4c0.2-0.3,0.2-0.7,0-1.1c-0.2-0.3-0.5-0.5-0.9-0.5L180.6,98.9L180.6,98.9
	L180.6,98.9z`})),_=e=>(0,s.createElement)(`svg`,Object.assign({xmlns:`http://www.w3.org/2000/svg`,viewBox:`0 0 100 100`,preserveAspectRatio:`xMidYMid`,width:84,height:84},e),(0,s.createElement)(`circle`,{cx:`50`,cy:`36.8101`,r:`10`,fill:`#3e7cff`},(0,s.createElement)(`animate`,{attributeName:`cy`,dur:`1s`,repeatCount:`indefinite`,calcMode:`spline`,keySplines:`0.45 0 0.9 0.55;0 0.45 0.55 0.9`,keyTimes:`0;0.5;1`,values:`23;77;23`})));function v(e){let t=0,n=0;if(e.getBoundingClientRect){let r=e.getBoundingClientRect(),i=document.documentElement;t=r.left+Math.max(i.scrollLeft,document.body.scrollLeft)-i.clientLeft,n=r.top+Math.max(i.scrollTop,document.body.scrollTop)-i.clientTop}else for(;e!==document.body;)t+=e.offsetLeft,n+=e.offsetTop,e=e.offsetParent;return{domX:t,domY:n}}function y(e,t){let n=t.relatedTarget;try{for(;n&&n!==e;)n=n.parentNode}catch(e){console.warn(e)}return n!==e}var b=(e,t,n)=>{let[r,i]=(0,s.useState)([]),a=(0,s.useCallback)(()=>{i([])},[i]),o=(0,s.useCallback)(e=>{let n=e.currentTarget,a=v(n),o=e.pageX||e.clientX,s=e.pageY||e.clientY,c=a.domX,l=a.domY,u=o-c,d=s-l,f=parseInt(u.toString()),p=parseInt(d.toString()),m=new Date,h=r.length;return i([...r,{key:m.getTime(),index:h+1,x:f,y:p}]),t.click&&t.click(f,p),e.cancelBubble=!0,e.preventDefault(),!1},[r,t]),c=(0,s.useCallback)(e=>(t.confirm&&t.confirm(r,()=>{a()}),e.cancelBubble=!0,e.preventDefault(),!1),[r,t,a]),l=(0,s.useCallback)(()=>r,[r]),u=(0,s.useCallback)(()=>{a(),n&&n()},[a,n]),d=(0,s.useCallback)(()=>{t.close&&t.close(),a()},[t,a]),f=(0,s.useCallback)(()=>{t.refresh&&t.refresh(),a()},[a]);return{setDots:i,getDots:l,clickEvent:o,confirmEvent:c,closeEvent:(0,s.useCallback)(e=>(d(),e.cancelBubble=!0,e.preventDefault(),!1),[d]),refreshEvent:(0,s.useCallback)(e=>(f(),e.cancelBubble=!0,e.preventDefault(),!1),[t,f]),resetData:a,clearData:u,close:d,refresh:f}},x=(0,s.forwardRef)((e,t)=>{let[n,r]=(0,s.useState)({...m(),...e.config||{}}),[i,a]=(0,s.useState)({...e.data||{}}),[o,l]=(0,s.useState)({...e.events||{}});(0,s.useEffect)(()=>{r({...n,...e.config||{}})},[e.config,r]),(0,s.useEffect)(()=>{a({...i,...e.data||{}})},[e.data,a]),(0,s.useEffect)(()=>{l({...o,...e.events||{}})},[e.events,l]);let u=b(i,o,()=>{a({...i,thumb:``,image:``})}),f=n.horizontalPadding||0,v=n.verticalPadding||0,y=(n.width||0)+f*2+(n.showTheme?2:0),x=(n.width||0)>0||(n.height||0)>0,S=i.image&&i.image.length>0&&i.thumb&&i.thumb.length>0;return(0,s.useImperativeHandle)(t,()=>({reset:u.resetData,clear:u.clearData,refresh:u.refresh,close:u.close})),s.createElement(`div`,{className:(0,c.default)(p.wrapper,n.showTheme?p.theme:``),style:{width:y+`px`,paddingLeft:f+`px`,paddingRight:f+`px`,paddingTop:v+`px`,paddingBottom:v+`px`,display:x?`block`:`none`}},s.createElement(`div`,{className:p.header},s.createElement(`span`,null,n.title),s.createElement(`img`,{className:i.thumb==``?p.hide:``,style:{width:n.thumbWidth+`px`,height:n.thumbHeight+`px`,display:S?`block`:`none`},src:i.thumb,alt:``})),s.createElement(`div`,{className:p.body,style:{width:n.width+`px`,height:n.height+`px`}},s.createElement(`div`,{className:p.loading},s.createElement(_,null)),s.createElement(`img`,{className:(0,c.default)(p.picture,i.image==``?p.hide:``),style:{width:n.width+`px`,height:n.height+`px`,display:S?`block`:`none`},src:i.image,alt:``,onClick:u.clickEvent}),s.createElement(`div`,{className:d.dots},u.getDots().map(e=>s.createElement(`div`,{className:`dot`,style:{width:n.dotSize+`px`,height:n.dotSize+`px`,borderRadius:n.dotSize+`px`,top:e.y-(n.dotSize||1)/2-1+`px`,left:e.x-(n.dotSize||1)/2-1+`px`},key:e.key+`-`+e.index},e.index)))),s.createElement(`div`,{className:p.footer},s.createElement(`div`,{className:(0,c.default)(p.iconBlock,d.iconBlock)},s.createElement(h,{width:n.iconSize,height:n.iconSize,onClick:u.closeEvent}),s.createElement(g,{width:n.iconSize,height:n.iconSize,onClick:u.refreshEvent})),s.createElement(`div`,{className:p.buttonBlock},s.createElement(`button`,{className:(0,c.default)(!S&&p.disabled),onClick:u.confirmEvent},n.buttonText))))}),S=s.memo(x),C=`/**
 * @Author Awen
 * @Date 2024/06/01
 * @Email wengaolng@gmail.com
 **/
.index-module_tile__8pkQD {
  position: absolute;
  z-index: 2;
  cursor: pointer;
}
.index-module_tile__8pkQD img {
  display: block;
  cursor: pointer;
  width: 100%;
  height: 100%;
}
`,w={tile:`index-module_tile__8pkQD`};l(C);var T=e=>(0,s.createElement)(`svg`,Object.assign({viewBox:`0 0 200 200`,xmlns:`http://www.w3.org/2000/svg`,width:20,height:20},e),(0,s.createElement)(`path`,{d:`M131.6,116.3c0,0-75.6,0-109.7,0c-9.1,0-16.2-7.4-16.2-16.2c0-9.1,7.4-16.2,16.2-16.2c28.7,0,109.7,0,109.7,0
	s-5.4-5.4-30.4-30.7c-6.4-6.4-6.4-16.7,0-23.1s16.7-6.4,23.1,0l58.4,58.4c6.4,6.4,6.4,16.7,0,23.1c0,0-32.9,32.9-57.9,57.9
	c-6.4,6.4-16.7,6.4-23.1,0c-6.4-6.4-6.4-16.7,0-23.1C121.8,126.2,131.6,116.3,131.6,116.3z`})),E=()=>({thumbX:0,thumbY:0,thumbWidth:0,thumbHeight:0,image:``,thumb:``}),D=()=>({width:300,height:220,thumbWidth:150,thumbHeight:40,verticalPadding:16,horizontalPadding:12,showTheme:!0,title:`请拖动滑块完成拼图`,iconSize:22,scope:!0}),O=(e,t,n,r,i,a,o,c,l)=>{let[u,d]=(0,s.useState)(0),[f,p]=(0,s.useState)(e.thumbX||0),[m,h]=(0,s.useState)(!1);(0,s.useEffect)(()=>{m||p(e.thumbX||0)},[e,p]);let g=(0,s.useCallback)(()=>{d(0),p(e.thumbX||0)},[d,p,e.thumbX]),_=(0,s.useCallback)(s=>{if(!y(c.current,s))return;let l=s.touches&&s.touches[0],u=o.current.offsetLeft,f=i.current.offsetWidth,m=f-o.current.offsetWidth,_=a.current.offsetWidth,v=a.current.offsetLeft,b=f-_,x=(f-(_+v))/m,S=!1,C=null,w=0,T=0;w=l?l.pageX-u:s.clientX-u;let E=n=>{S=!0;let r=n.touches&&n.touches[0],i=0;i=r?r.pageX-w:n.clientX-w;let a=v+i*x;if(i>=m){d(m),T=b,p(T);return}if(i<=0){d(0),T=v,p(T);return}d(i),T=T=a,p(T),t.move&&t.move(T,e.thumbY||0),n.cancelBubble=!0,n.preventDefault()},D=n=>{y(c.current,n)&&(P(),S&&(S=!1,!(T<0)&&(t.confirm&&t.confirm({x:parseInt(T.toString()),y:e.thumbY||0},()=>{g()}),n.cancelBubble=!0,n.preventDefault())))},O=e=>{C=e},k=()=>{C=null},A=e=>{C&&(D(C),P())},j=n.scope,M=j?r.current:c.current,N=j?r.current:document.body,P=()=>{N.removeEventListener(`mousemove`,E,!1),N.removeEventListener(`touchmove`,E,{passive:!1}),M.removeEventListener(`mouseup`,D,!1),M.removeEventListener(`mouseenter`,k,!1),M.removeEventListener(`mouseleave`,O,!1),M.removeEventListener(`touchend`,D,!1),N.removeEventListener(`mouseleave`,D,!1),N.removeEventListener(`mouseup`,A,!1),h(!1)};h(!0),N.addEventListener(`mousemove`,E,!1),N.addEventListener(`touchmove`,E,{passive:!1}),M.addEventListener(`mouseup`,D,!1),M.addEventListener(`mouseenter`,k,!1),M.addEventListener(`mouseleave`,O,!1),M.addEventListener(`touchend`,D,!1),N.addEventListener(`mouseleave`,D,!1),N.addEventListener(`mouseup`,A,!1)},[r,o,i,n,e,a,c,t,g]),v=(0,s.useCallback)(()=>{g(),l&&l()},[g,l]),b=(0,s.useCallback)(()=>{t.close&&t.close(),g()},[t,g]),x=(0,s.useCallback)(()=>{t.refresh&&t.refresh(),g()},[t,g]),S=(0,s.useCallback)(e=>(b(),e.cancelBubble=!0,e.preventDefault(),!1),[b]),C=(0,s.useCallback)(e=>(x(),e.cancelBubble=!0,e.preventDefault(),!1),[x]),w=(0,s.useCallback)(()=>({x:f,y:e.thumbY||0}),[e,f]);return{getState:(0,s.useCallback)(()=>({dragLeft:u,thumbLeft:f}),[f,u]),getPoint:w,dragEvent:_,closeEvent:S,refreshEvent:C,resetData:g,clearData:v,close:b,refresh:x}},k=(0,s.forwardRef)((e,t)=>{let[n,r]=(0,s.useState)({...D(),...e.config||{}}),[i,a]=(0,s.useState)({...E(),...e.data||{}}),[o,l]=(0,s.useState)({...e.events||{}});(0,s.useEffect)(()=>{r({...n,...e.config||{}})},[e.config,r]),(0,s.useEffect)(()=>{a({...i,...e.data||{}})},[e.data,a]),(0,s.useEffect)(()=>{l({...o,...e.events||{}})},[e.events,l]);let u=(0,s.useRef)(null),d=(0,s.useRef)(null),f=(0,s.useRef)(null),m=(0,s.useRef)(null),v=(0,s.useRef)(null),y=O(i,o,n,u,f,v,m,d,()=>{a({...i,...E()})}),b=n.horizontalPadding||0,x=n.verticalPadding||0,S=(n.width||0)+b*2+(n.showTheme?2:0),C=(n.width||0)>0||(n.height||0)>0,k=i.image&&i.image.length>0&&i.thumb&&i.thumb.length>0;return(0,s.useImperativeHandle)(t,()=>({reset:y.resetData,clear:y.clearData,refresh:y.refresh,close:y.close})),(0,s.useEffect)(()=>{let e=e=>e.preventDefault();return m.current&&m.current.addEventListener(`dragstart`,e),()=>{m.current&&m.current.removeEventListener(`dragstart`,e)}},[m]),s.createElement(`div`,{className:(0,c.default)(p.wrapper,n.showTheme?p.theme:``),style:{width:S+`px`,paddingLeft:b+`px`,paddingRight:b+`px`,paddingTop:x+`px`,paddingBottom:x+`px`,display:C?`block`:`none`},ref:u},s.createElement(`div`,{className:p.header},s.createElement(`span`,null,n.title),s.createElement(`div`,{className:p.iconBlock},s.createElement(h,{width:n.iconSize,height:n.iconSize,onClick:y.closeEvent}),s.createElement(g,{width:n.iconSize,height:n.iconSize,onClick:y.refreshEvent}))),s.createElement(`div`,{className:p.body,ref:f,style:{width:n.width+`px`,height:n.height+`px`}},s.createElement(`div`,{className:p.loading},s.createElement(_,null)),s.createElement(`img`,{className:(0,c.default)(p.picture,i.image==``?p.hide:``),style:{width:n.width+`px`,height:n.height+`px`,display:k?`block`:`none`},src:i.image,alt:``}),s.createElement(`div`,{className:w.tile,ref:v,style:{width:(i.thumbWidth||0)+`px`,height:(i.thumbHeight||0)+`px`,top:(i.thumbY||0)+`px`,left:y.getState().thumbLeft+`px`}},s.createElement(`img`,{className:i.thumb==``?p.hide:``,style:{display:k?`block`:`none`},src:i.thumb,alt:``}))),s.createElement(`div`,{className:p.footer},s.createElement(`div`,{className:p.dragSlideBar,ref:d},s.createElement(`div`,{className:p.dragLine}),s.createElement(`div`,{className:(0,c.default)(p.dragBlock,!k&&p.disabled),ref:m,onMouseDown:y.dragEvent,style:{left:y.getState().dragLeft+`px`}},s.createElement(`div`,{className:p.dragBlockInline,onTouchStart:y.dragEvent},s.createElement(T,null))))))}),A=s.memo(k),j=`/**
 * @Author Awen
 * @Date 2024/06/01
 * @Email wengaolng@gmail.com
 **/
.index-module_header__jVeEs {
  text-align: center;
}
.index-module_tile__VR9Ut {
  position: absolute;
  z-index: 2;
  cursor: pointer;
  -webkit-touch-callout: none;
  -webkit-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
.index-module_tile__VR9Ut img {
  display: block;
  cursor: pointer;
  width: 100%;
  height: 100%;
}
`,M={header:`index-module_header__jVeEs`,tile:`index-module_tile__VR9Ut`};l(j);var N=()=>({width:300,height:220,verticalPadding:16,horizontalPadding:12,showTheme:!0,title:`请拖动滑块完成拼图`,iconSize:22,scope:!0}),P=(e,t,n,r,i,a,o)=>{let[c,l]=(0,s.useState)({x:e.thumbX||0,y:e.thumbY||0}),[u,d]=(0,s.useState)(!1);(0,s.useEffect)(()=>{u||l({x:e.thumbX||0,y:e.thumbY||0})},[e,l]);let f=(0,s.useCallback)(()=>{l({x:e.thumbX||0,y:e.thumbY||0})},[e.thumbX,e.thumbY,l]),p=(0,s.useCallback)(e=>{if(!y(i.current,e))return;let o=e.touches&&e.touches[0],s=a.current.offsetLeft,c=a.current.offsetTop,u=i.current.offsetWidth,p=i.current.offsetHeight,m=a.current.offsetWidth,h=a.current.offsetHeight,g=u-m,_=p-h,v=!1,b=null,x=0,S=0,C=0,w=0;o?(x=o.pageX-s,S=o.pageY-c):(x=e.clientX-s,S=e.clientY-c);let T=e=>{v=!0;let n=e.touches&&e.touches[0],r=0,i=0;n?(r=n.pageX-x,i=n.pageY-S):(r=e.clientX-x,i=e.clientY-S),r<=0&&(r=0),i<=0&&(i=0),r>=g&&(r=g),i>=_&&(i=_),l({x:r,y:i}),C=r,w=i,t.move&&t.move(r,i),e.cancelBubble=!0,e.preventDefault()},E=e=>{y(i.current,e)&&(N(),v&&(v=!1,!(C<0||w<0)&&(t.confirm&&t.confirm({x:C,y:w},()=>{f()}),e.cancelBubble=!0,e.preventDefault())))},D=e=>{b=e},O=()=>{b=null},k=e=>{b&&(E(b),N())},A=n.scope,j=A?r.current:i.current,M=A?r.current:document.body,N=()=>{M.removeEventListener(`mousemove`,T,!1),M.removeEventListener(`touchmove`,T,{passive:!1}),j.removeEventListener(`mouseup`,E,!1),j.removeEventListener(`mouseenter`,O,!1),j.removeEventListener(`mouseleave`,D,!1),j.removeEventListener(`touchend`,E,!1),M.removeEventListener(`mouseleave`,E,!1),M.removeEventListener(`mouseup`,k,!1),d(!1)};d(!0),M.addEventListener(`mousemove`,T,!1),M.addEventListener(`touchmove`,T,{passive:!1}),j.addEventListener(`mouseup`,E,!1),j.addEventListener(`mouseenter`,O,!1),j.addEventListener(`mouseleave`,D,!1),j.addEventListener(`touchend`,E,!1),M.addEventListener(`mouseleave`,E,!1),M.addEventListener(`mouseup`,k,!1)},[r,i,a,n,t,d,f]),m=(0,s.useCallback)(()=>{f(),o&&o()},[f,o]),h=(0,s.useCallback)(()=>{t.close&&t.close(),f()},[t,f]),g=(0,s.useCallback)(()=>{t.refresh&&t.refresh(),f()},[t,f]);return{thumbPoint:c,dragEvent:p,closeEvent:(0,s.useCallback)(e=>(h(),e.cancelBubble=!0,e.preventDefault(),!1),[h]),refreshEvent:(0,s.useCallback)(e=>(g(),e.cancelBubble=!0,e.preventDefault(),!1),[g]),resetData:f,clearData:m,close:h,refresh:g}},F=()=>({thumbX:0,thumbY:0,thumbWidth:0,thumbHeight:0,image:``,thumb:``}),I=(0,s.forwardRef)((e,t)=>{let[n,r]=(0,s.useState)({...N(),...e.config||{}}),[i,a]=(0,s.useState)({...F(),...e.data||{}}),[o,l]=(0,s.useState)({...e.events||{}});(0,s.useEffect)(()=>{r({...n,...e.config||{}})},[e.config,r]),(0,s.useEffect)(()=>{a({...i,...e.data||{}})},[e.data,a]),(0,s.useEffect)(()=>{l({...o,...e.events||{}})},[e.events,l]);let u=(0,s.useRef)(null),d=(0,s.useRef)(null),f=(0,s.useRef)(null),m=P(i,o,n,u,d,f,()=>{a({...i,...F()})}),v=n.horizontalPadding||0,y=n.verticalPadding||0,b=(n.width||0)+v*2+(n.showTheme?2:0),x=(n.width||0)>0||(n.height||0)>0,S=i.image&&i.image.length>0&&i.thumb&&i.thumb.length>0;return(0,s.useImperativeHandle)(t,()=>({reset:m.resetData,clear:m.clearData,refresh:m.refresh,close:m.close})),(0,s.useEffect)(()=>{let e=e=>e.preventDefault();return f.current&&f.current.addEventListener(`dragstart`,e),()=>{f.current&&f.current.removeEventListener(`dragstart`,e)}},[f]),s.createElement(`div`,{className:(0,c.default)(p.wrapper,M.wrapper,n.showTheme?p.theme:``),style:{width:b+`px`,paddingLeft:v+`px`,paddingRight:v+`px`,paddingTop:y+`px`,paddingBottom:y+`px`,display:x?`block`:`none`},ref:u},s.createElement(`div`,{className:(0,c.default)(p.header,M.header)},s.createElement(`span`,null,n.title)),s.createElement(`div`,{className:p.body,ref:d,style:{width:n.width+`px`,height:n.height+`px`}},s.createElement(`div`,{className:p.loading},s.createElement(_,null)),s.createElement(`img`,{className:(0,c.default)(p.picture,i.image==``?p.hide:``),src:i.image,style:{width:n.width+`px`,height:n.height+`px`,display:S?`block`:`none`},alt:``}),s.createElement(`div`,{className:M.tile,ref:f,style:{width:(i.thumbWidth||0)+`px`,height:(i.thumbHeight||0)+`px`,top:m.thumbPoint.y+`px`,left:m.thumbPoint.x+`px`},onMouseDown:m.dragEvent,onTouchStart:m.dragEvent},s.createElement(`img`,{className:i.thumb==``?p.hide:``,style:{display:S?`block`:`none`},src:i.thumb,alt:``}))),s.createElement(`div`,{className:p.footer},s.createElement(`div`,{className:p.iconBlock},s.createElement(h,{width:n.iconSize,height:n.iconSize,onClick:m.closeEvent}),s.createElement(g,{width:n.iconSize,height:n.iconSize,onClick:m.refreshEvent}))))}),L=s.memo(I),R=`/**
 * @Author Awen
 * @Date 2024/06/01
 * @Email wengaolng@gmail.com
 **/
.index-module_body__5eTaZ {
  background: transparent !important;
  display: flex;
  display: -webkit-flex;
  justify-content: center;
  align-items: center;
  margin: 10px auto 0;
}
.index-module_bodyInner__Lb3mp {
  border-radius: 100%;
}
.index-module_picture__M-qbX {
  position: relative;
  max-width: 100%;
  max-height: 100%;
  z-index: 2;
  border-radius: 100%;
  overflow: hidden;
  display: flex;
  justify-content: center;
  align-items: center;
}
.index-module_picture__M-qbX img {
  max-width: 100%;
  max-height: 100%;
}
.index-module_round__zaOPS {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 100%;
  z-index: 2;
  border: 6px solid #e0e0e0;
  border-color: var(--go-captcha-theme-round-color);
}
.index-module_thumb__jChIh {
  position: absolute;
  z-index: 2;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: center;
  align-items: center;
}
.index-module_thumb__jChIh img {
  max-width: 100%;
  max-height: 100%;
}
.index-module_thumbBlock__u3U1X {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}
`,z={body:`index-module_body__5eTaZ`,bodyInner:`index-module_bodyInner__Lb3mp`,picture:`index-module_picture__M-qbX`,round:`index-module_round__zaOPS`,thumb:`index-module_thumb__jChIh`,thumbBlock:`index-module_thumbBlock__u3U1X`};l(R);var B=()=>({width:300,height:220,size:220,verticalPadding:16,horizontalPadding:12,showTheme:!0,title:`请拖动滑块完成拼图`,iconSize:22,scope:!0}),V=(e,t,n,r,i,a,o)=>{let[c,l]=(0,s.useState)(0),[u,d]=(0,s.useState)(e.angle||0),[f,p]=(0,s.useState)(!1);(0,s.useEffect)(()=>{f||d(e.angle||0)},[e,d]);let m=(0,s.useCallback)(()=>{l(0),d(e.angle||0)},[e.angle,l,d]),h=(0,s.useCallback)(o=>{if(!y(a.current,o))return;let s=o.touches&&o.touches[0],c=i.current.offsetLeft,u=a.current.offsetWidth-i.current.offsetWidth,f=(360-e.angle)/u,h=0,g=!1,_=null,v=0,b=0;v=s?s.pageX-c:o.clientX-c;let x=n=>{g=!0;let r=n.touches&&n.touches[0],i=0;if(i=r?r.pageX-v:n.clientX-v,h=e.angle+i*f,i>=u){l(u),b=360,d(b);return}if(i<=0){l(0),b=e.angle,d(b);return}l(i),b=h,d(h),t.rotate&&t.rotate(h),n.cancelBubble=!0,n.preventDefault()},S=e=>{y(a.current,e)&&(k(),g&&(g=!1,!(b<0)&&(t.confirm&&t.confirm(parseInt(b.toString()),()=>{m()}),e.cancelBubble=!0,e.preventDefault())))},C=e=>{_=e},w=()=>{_=null},T=e=>{_&&(S(_),k())},E=n.scope,D=E?r.current:a.current,O=E?r.current:document.body,k=()=>{O.removeEventListener(`mousemove`,x,!1),O.removeEventListener(`touchmove`,x,{passive:!1}),D.removeEventListener(`mouseup`,S,!1),D.removeEventListener(`mouseenter`,w,!1),D.removeEventListener(`mouseleave`,C,!1),D.removeEventListener(`touchend`,S,!1),O.removeEventListener(`mouseleave`,S,!1),O.removeEventListener(`mouseup`,T,!1),p(!1)};p(!0),O.addEventListener(`mousemove`,x,!1),O.addEventListener(`touchmove`,x,{passive:!1}),D.addEventListener(`mouseup`,S,!1),D.addEventListener(`mouseenter`,w,!1),D.addEventListener(`mouseleave`,C,!1),D.addEventListener(`touchend`,S,!1),O.addEventListener(`mouseleave`,S,!1),O.addEventListener(`mouseup`,T,!1)},[r,i,a,n,e,t,m]),g=(0,s.useCallback)(()=>{m(),o&&o()},[m,o]),_=(0,s.useCallback)(()=>{t.close&&t.close(),m()},[t,m]),v=(0,s.useCallback)(()=>{t.refresh&&t.refresh(),m()},[t,m]),b=(0,s.useCallback)(e=>(_(),e.cancelBubble=!0,e.preventDefault(),!1),[_]),x=(0,s.useCallback)(e=>(v(),e.cancelBubble=!0,e.preventDefault(),!1),[v]);return{getState:(0,s.useCallback)(()=>({dragLeft:c,thumbAngle:u}),[u,c]),thumbAngle:u,dragEvent:h,closeEvent:b,refreshEvent:x,resetData:m,clearData:g,close:_,refresh:v}},H=()=>({angle:0,image:``,thumb:``,thumbSize:0}),U=(0,s.forwardRef)((e,t)=>{let[n,r]=(0,s.useState)({...B(),...e.config||{}}),[i,a]=(0,s.useState)({...H(),...e.data||{}}),[o,l]=(0,s.useState)({...e.events||{}});(0,s.useEffect)(()=>{r({...n,...e.config||{}})},[e.config,r]),(0,s.useEffect)(()=>{a({...i,...e.data||{}})},[e.data,a]),(0,s.useEffect)(()=>{l({...o,...e.events||{}})},[e.events,l]);let u=(0,s.useRef)(null),d=(0,s.useRef)(null),f=(0,s.useRef)(null),m=V(i,o,n,u,f,d,()=>{a({...i,...H()})}),v=n.horizontalPadding||0,y=n.verticalPadding||0,b=(n.width||0)+v*2+(n.showTheme?2:0),x=(n.size||0)>0?n.size:B().size,S=(n.width||0)>0||(n.height||0)>0,C=i.image&&i.image.length>0&&i.thumb&&i.thumb.length>0;return(0,s.useImperativeHandle)(t,()=>({reset:m.resetData,clear:m.clearData,refresh:m.refresh,close:m.close})),(0,s.useEffect)(()=>{let e=e=>e.preventDefault();return f.current&&f.current.addEventListener(`dragstart`,e),()=>{f.current&&f.current.removeEventListener(`dragstart`,e)}},[f]),s.createElement(`div`,{className:(0,c.default)(p.wrapper,z.wrapper,n.showTheme?p.theme:``),style:{width:b+`px`,paddingLeft:v+`px`,paddingRight:v+`px`,paddingTop:y+`px`,paddingBottom:y+`px`,display:S?`block`:`none`},ref:u},s.createElement(`div`,{className:p.header},s.createElement(`span`,null,n.title),s.createElement(`div`,{className:p.iconBlock},s.createElement(h,{width:n.iconSize,height:n.iconSize,onClick:m.closeEvent}),s.createElement(g,{width:n.iconSize,height:n.iconSize,onClick:m.refreshEvent}))),s.createElement(`div`,{className:(0,c.default)(p.body,z.body),style:{width:n.width+`px`,height:n.height+`px`}},s.createElement(`div`,{className:(0,c.default)(z.bodyInner,p.bodyInner),style:{width:x+`px`,height:x+`px`}},s.createElement(`div`,{className:p.loading},s.createElement(_,null)),s.createElement(`div`,{className:z.picture,style:{width:n.size+`px`,height:n.size+`px`}},s.createElement(`img`,{className:i.image==``?p.hide:``,src:i.image,style:{display:C?`block`:`none`},alt:``}),s.createElement(`div`,{className:z.round})),s.createElement(`div`,{className:z.thumb},s.createElement(`div`,{className:z.thumbBlock,style:{transform:`rotate(`+m.getState().thumbAngle+`deg)`,...i.thumbSize>0?{width:i.thumbSize+`px`,height:i.thumbSize+`px`}:{}}},s.createElement(`img`,{className:i.thumb==``?p.hide:``,src:i.thumb,style:{visibility:C?`visible`:`hidden`},alt:``}))))),s.createElement(`div`,{className:p.footer},s.createElement(`div`,{className:p.dragSlideBar,ref:d},s.createElement(`div`,{className:p.dragLine}),s.createElement(`div`,{className:(0,c.default)(p.dragBlock,!C&&p.disabled),ref:f,onMouseDown:m.dragEvent,style:{left:m.getState().dragLeft+`px`}},s.createElement(`div`,{className:p.dragBlockInline,onTouchStart:m.dragEvent},s.createElement(T,null))))))}),W=s.memo(U),G=()=>({width:330,height:44,verticalPadding:12,horizontalPadding:16}),K=`/**
 * @Author Awen
 * @Date 2024/06/01
 * @Email wengaolng@gmail.com
 **/
.index-module_btnBlock__L96Vx {
  position: relative;
  box-sizing: border-box;
  display: block;
  font-size: 13px;
  -webkit-border-radius: 5px;
  -moz-border-radius: 5px;
  letter-spacing: 1px;
  border-radius: 5px;
  line-height: 1;
  white-space: nowrap;
  -webkit-appearance: none;
  outline: none;
  margin: 0;
  transition: 0.1s;
  font-weight: 500;
  -moz-user-select: none;
  -webkit-user-select: none;
  display: flex;
  align-items: center;
  justify-content: center;
  justify-items: center;
  box-shadow: 0 0 20px rgba(62, 124, 255, 0.1);
  -webkit-box-shadow: 0 0 20px rgba(62, 124, 255, 0.1);
  -moz-box-shadow: 0 0 20px rgba(62, 124, 255, 0.1);
}
.index-module_btnBlock__L96Vx span {
  padding-left: 8px;
}
.index-module_disabled__U5sNo {
  pointer-events: none;
}
.index-module_default__r2sQq {
  color: var(--go-captcha-theme-default-color);
  border: 1px solid #50a1ff;
  border-color: var(--go-captcha-theme-default-border-color);
  background-color: var(--go-captcha-theme-default-bg-color);
  cursor: pointer;
}
.index-module_default__r2sQq:hover {
  background-color: var(--go-captcha-theme-default-hover-color) !important;
}
.index-module_error__mCm6a {
  cursor: pointer;
  color: var(--go-captcha-theme-error-color);
  background-color: var(--go-captcha-theme-error-bg-color);
  border: 1px solid #ff5a34;
  border-color: var(--go-captcha-theme-error-border-color);
}
.index-module_warn__CT1sW {
  cursor: pointer;
  color: var(--go-captcha-theme-warn-color);
  background-color: var(--go-captcha-theme-warn-bg-color);
  border: 1px solid #ffbe09;
  border-color: var(--go-captcha-theme-warn-border-color);
}
.index-module_success__61kOU {
  color: var(--go-captcha-theme-success-color);
  background-color: var(--go-captcha-theme-success-bg-color);
  border: 1px solid #8bc640;
  border-color: var(--go-captcha-theme-success-border-color);
  pointer-events: none;
}
.index-module_ripple__KF4IK {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  justify-items: center;
}
.index-module_ripple__KF4IK svg {
  position: relative;
  z-index: 2;
}
.index-module_ripple__KF4IK > * {
  z-index: 2;
}
.index-module_ripple__KF4IK::after {
  background-color: var(--go-captcha-theme-default-border-color);
  border-radius: 50px;
  content: "";
  display: block;
  width: 21px;
  height: 21px;
  opacity: 0;
  position: absolute;
  top: 50%;
  left: 50%;
  margin-top: -11px;
  margin-left: -11px;
  z-index: 1;
  animation: index-module_ripple__KF4IK 1.3s infinite;
  -moz-animation: index-module_ripple__KF4IK 1.3s infinite;
  -webkit-animation: index-module_ripple__KF4IK 1.3s infinite;
  animation-delay: 2s;
  -moz-animation-delay: 2s;
  -webkit-animation-delay: 2s;
}
@keyframes index-module_ripple__KF4IK {
  0% {
    opacity: 0;
  }
  5% {
    opacity: 0.05;
  }
  20% {
    opacity: 0.35;
  }
  65% {
    opacity: 0.01;
  }
  100% {
    transform: scaleX(2) scaleY(2);
    opacity: 0;
  }
}
`,q={btnBlock:`index-module_btnBlock__L96Vx`,disabled:`index-module_disabled__U5sNo`,default:`index-module_default__r2sQq`,error:`index-module_error__mCm6a`,warn:`index-module_warn__CT1sW`,success:`index-module_success__61kOU`,ripple:`index-module_ripple__KF4IK`};l(K);var J=e=>(0,s.createElement)(`svg`,Object.assign({viewBox:`0 0 200 200`,xmlns:`http://www.w3.org/2000/svg`,width:20,height:20},e),(0,s.createElement)(`circle`,{fill:`#3E7CFF`,cx:`100`,cy:`100`,r:`96.3`}),(0,s.createElement)(`path`,{fill:`#FFFFFF`,d:`M140.8,64.4l-39.6-11.9h-2.4L59.2,64.4c-1.6,0.8-2.8,2.4-2.8,4v24.1c0,25.3,15.8,45.9,42.3,54.6
	c0.4,0,0.8,0.4,1.2,0.4c0.4,0,0.8,0,1.2-0.4c26.5-8.7,42.3-28.9,42.3-54.6V68.3C143.5,66.8,142.3,65.2,140.8,64.4z`})),Y=e=>(0,s.createElement)(`svg`,Object.assign({viewBox:`0 0 200 200`,xmlns:`http://www.w3.org/2000/svg`,width:20,height:20},e),(0,s.createElement)(`path`,{fill:`#ED4630`,d:`M184,26.6L102.4,2.1h-4.9L16,26.6c-3.3,1.6-5.7,4.9-5.7,8.2v49.8c0,52.2,32.6,94.7,87.3,112.6
	c0.8,0,1.6,0.8,2.4,0.8s1.6,0,2.4-0.8c54.7-18,87.3-59.6,87.3-112.6V34.7C189.8,31.5,187.3,28.2,184,26.6z M134.5,123.1
	c3.1,3.1,3.1,8.2,0,11.3c-1.6,1.6-3.6,2.3-5.7,2.3s-4.1-0.8-5.7-2.3L100,111.3l-23.1,23.1c-1.6,1.6-3.6,2.3-5.7,2.3
	c-2,0-4.1-0.8-5.7-2.3c-3.1-3.1-3.1-8.2,0-11.3L88.7,100L65.5,76.9c-3.1-3.1-3.1-8.2,0-11.3c3.1-3.1,8.2-3.1,11.3,0L100,88.7
	l23.1-23.1c3.1-3.1,8.2-3.1,11.3,0c3.1,3.1,3.1,8.2,0,11.3L111.3,100L134.5,123.1z`})),X=e=>(0,s.createElement)(`svg`,Object.assign({viewBox:`0 0 200 200`,xmlns:`http://www.w3.org/2000/svg`,width:20,height:20},e),(0,s.createElement)(`path`,{fill:`#FFA000`,d:`M184,26.6L102.4,2.1h-4.9L16,26.6c-3.3,1.6-5.7,4.9-5.7,8.2v49.8c0,52.2,32.6,94.7,87.3,112.6
	c0.8,0,1.6,0.8,2.4,0.8s1.6,0,2.4-0.8c54.7-18,87.3-59.6,87.3-112.6V34.7C189.8,31.5,187.3,28.2,184,26.6z M107.3,109.1
	c-0.5,5.4-3.9,7.9-7.3,7.9c-2.5,0,0,0,0,0c-3.2-0.6-5.7-2-6.8-7.4l-4.4-50.9c0-5.1,6.2-9.7,11.5-9.7c5.3,0,11,4.7,11,9.9
	L107.3,109.1z M109.3,133.3c0,5.1-4.2,9.3-9.3,9.3c-5.1,0-9.3-4.2-9.3-9.3c0-5.1,4.2-9.3,9.3-9.3C105.1,124,109.3,128.1,109.3,133.3
	z`})),Z=e=>(0,s.createElement)(`svg`,Object.assign({viewBox:`0 0 200 200`,xmlns:`http://www.w3.org/2000/svg`,width:20,height:20},e),(0,s.createElement)(`path`,{fill:`#5EAA2F`,d:`M183.3,27.2L102.4,2.9h-4.9L16.7,27.2C13.4,28.8,11,32,11,35.3v49.4c0,51.8,32.4,93.9,86.6,111.7
	c0.8,0,1.6,0.8,2.4,0.8c0.8,0,1.6,0,2.4-0.8c54.2-17.8,86.6-59.1,86.6-111.7V35.3C189,32,186.6,28.8,183.3,27.2z M146.1,81.4
	l-48.5,48.5c-1.6,1.6-3.2,2.4-5.7,2.4c-2.4,0-4-0.8-5.7-2.4L62,105.7c-3.2-3.2-3.2-8.1,0-11.3c3.2-3.2,8.1-3.2,11.3,0l18.6,18.6
	l42.9-42.9c3.2-3.2,8.1-3.2,11.3,0C149.4,73.3,149.4,78.2,146.1,81.4L146.1,81.4z`})),Q={Click:S,Slide:A,SlideRegion:L,Rotate:W,Button:s.memo(e=>{let[t,n]=(0,s.useState)({...G(),...e.config||{}});(0,s.useEffect)(()=>{n({...t,...e.config||{}})},[e.config]);let r=e.type||`default`,i=s.createElement(J,null),a=q.default;return r==`warn`?(i=s.createElement(X,null),a=q.warn):r==`error`?(i=s.createElement(Y,null),a=q.error):r==`success`&&(i=s.createElement(Z,null),a=q.success),s.createElement(`div`,{className:(0,c.default)(q.btnBlock,a,e.disabled?q.disabled:``),style:{width:t.width+`px`,height:t.height+`px`,paddingLeft:t.verticalPadding+`px`,paddingRight:t.verticalPadding+`px`,paddingTop:t.verticalPadding+`px`,paddingBottom:t.verticalPadding+`px`},onClick:e.clickEvent},r==`default`?s.createElement(`div`,{className:q.ripple},i):i,s.createElement(`span`,null,e.title||`点击按键进行验证`))})},$=r(),ee=(0,s.forwardRef)(function({onVerified:e},t){let{t:n}=a(),[r,o]=(0,s.useState)(null),[c,l]=(0,s.useState)(!1),[u,d]=(0,s.useState)(!1),[f,p]=(0,s.useState)(null),m=(0,s.useRef)(null),h=(0,s.useRef)(null),g=(0,s.useRef)(null),_=(0,s.useCallback)(async()=>{l(!0),p(null);try{o(await i(`/api/captcha/request_angle`))}catch{o(null),p(n(`captcha.loadFailed`))}finally{l(!1)}},[n]);return(0,s.useImperativeHandle)(t,()=>({open:async()=>{d(!0),await _()},reset:()=>{m.current&&(m.current.value=``),h.current&&(h.current.value=``),g.current&&(g.current.value=`2`),o(null),d(!1)},hasCaptcha:()=>!!(m.current?.value&&h.current?.value),getCaptcha:()=>({captchaId:m.current?.value||``,captchaCode:h.current?.value||``,captchaProtocol:Number(g.current?.value)||2})}),[_]),(0,$.jsxs)($.Fragment,{children:[(0,$.jsx)(`input`,{ref:m,type:`hidden`,name:`captchaId`}),(0,$.jsx)(`input`,{ref:h,type:`hidden`,name:`captchaCode`}),(0,$.jsx)(`input`,{ref:g,type:`hidden`,name:`captchaProtocol`,value:`2`,readOnly:!0}),u?(0,$.jsx)(`div`,{className:`fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4`,children:(0,$.jsx)(`div`,{className:`w-auto max-w-none overflow-hidden rounded-lg bg-white p-0 shadow-lg`,children:r?(0,$.jsx)(Q.Rotate,{config:{title:n(`captcha.title`)},data:{image:r.imageBase64,thumb:r.thumbBase64,thumbSize:r.thumbSize,angle:0},events:{refresh:()=>{_()},close:()=>d(!1),confirm:t=>{m.current&&(m.current.value=r.id),h.current&&(h.current.value=String(t)),d(!1),e()}}}):(0,$.jsx)(`div`,{className:`flex min-h-60 min-w-80 items-center justify-center p-6 text-sm text-muted-foreground`,children:c?n(`captcha.loading`):f||n(`captcha.loadFailed`)})})}):null]})});export{ee as t};