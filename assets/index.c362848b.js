import{_ as I,g as B,B as k,b as F,y as P}from"./index.1b6e845e.js";/* empty css              *//* empty css               *//* empty css               *//* empty css                *//* empty css               *//* empty css                *//* empty css              */import{e as T,r as h,c as D,o as U,b4 as j,bE as E,ab as J,aV as N,b8 as G,bD as H,bF as A,C as V,F as c,aH as e,aG as a,u as t,aD as L,B as _,aJ as M,b5 as O,bG as q}from"./arco.ef2238db.js";import K from"./TopicList.a5dc0dfc.js";import"./chart.ba8df0e9.js";import"./vue.f27dce51.js";/* empty css               *//* empty css               *//* empty css               */const Q={class:"container"},R={class:"container-header"},W={class:"container-main"},X={key:0,class:"topic-container"},Y={class:"container-footer"},Z={__name:"index",setup($){const b=B(),g=T(!1),o=h({limit:20,page:1}),n=h({page:{page:1,limit:20,total:0},results:[]}),p=D(()=>({total:n.page.total,current:n.page.page,pageSize:n.page.limit,showTotal:!0,showJumper:!0,showPageSize:!0}));U(()=>{k()});const i=async()=>{g.value=!0;try{const r=await F.postForm("/api/admin/topic/list",P(o));n.page=r.page,n.results=r.results}finally{g.value=!1}};i();const v=r=>{o.page=r,i()},w=r=>{o.limit=r,i()};return(r,l)=>{const d=j,u=O,m=q,f=E,y=J,z=N,x=G,C=H,S=A;return _(),V("div",Q,[c("div",R,[e(x,{model:t(o),layout:"inline",size:t(b).table.size},{default:a(()=>[e(u,null,{default:a(()=>[e(d,{modelValue:t(o).id,"onUpdate:modelValue":l[0]||(l[0]=s=>t(o).id=s),placeholder:"ID"},null,8,["modelValue"])]),_:1}),e(u,null,{default:a(()=>[e(d,{modelValue:t(o).userId,"onUpdate:modelValue":l[1]||(l[1]=s=>t(o).userId=s),placeholder:"\u7528\u6237ID"},null,8,["modelValue"])]),_:1}),e(u,null,{default:a(()=>[e(f,{modelValue:t(o).status,"onUpdate:modelValue":l[2]||(l[2]=s=>t(o).status=s),placeholder:"\u72B6\u6001","allow-clear":"",onChange:i},{default:a(()=>[e(m,{value:0,label:"\u6B63\u5E38"}),e(m,{value:1,label:"\u5220\u9664"}),e(m,{value:2,label:"\u5F85\u5BA1\u6838"})]),_:1},8,["modelValue"])]),_:1}),e(u,null,{default:a(()=>[e(f,{modelValue:t(o).recommend,"onUpdate:modelValue":l[3]||(l[3]=s=>t(o).recommend=s),placeholder:"\u662F\u5426\u63A8\u8350","allow-clear":"",onChange:i},{default:a(()=>[e(m,{value:1,label:"\u63A8\u8350"}),e(m,{value:0,label:"\u672A\u63A8\u8350"})]),_:1},8,["modelValue"])]),_:1}),e(u,null,{default:a(()=>[e(d,{modelValue:t(o).title,"onUpdate:modelValue":l[4]||(l[4]=s=>t(o).title=s),placeholder:"\u6807\u9898"},null,8,["modelValue"])]),_:1}),e(u,null,{default:a(()=>[e(z,{type:"primary","html-type":"submit",onClick:i},{icon:a(()=>[e(y)]),default:a(()=>[M(" \u67E5\u8BE2 ")]),_:1})]),_:1})]),_:1},8,["model","size"])]),c("div",W,[t(n)&&t(n).results?(_(),V("div",X,[e(K,{results:t(n).results,onChange:i},null,8,["results"])])):(_(),L(C,{key:1}))]),c("div",Y,[e(S,{style:{margin:"10px"},total:t(p).total,current:t(p).current,"page-size":t(p).pageSize,"show-total":t(p).showTotal,"show-jumper":t(p).showJumper,"show-page-size":t(p).showPageSize,onChange:v,onPageSizeChange:w},null,8,["total","current","page-size","show-total","show-jumper","show-page-size"])])])}}},ge=I(Z,[["__scopeId","data-v-0d4e88af"]]);export{ge as default};
