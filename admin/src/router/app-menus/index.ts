import { appRoutes, appExternalRoutes } from '../routes';

const mixinRoutes = [...appRoutes, ...appExternalRoutes];

// 原代码
// const appClientMenus = mixinRoutes.map((el) => {
//   const { name, path, meta, redirect, children } = el;
//   return {
//     name,
//     path,
//     meta,
//     redirect,
//     children,
//   };
// });

// 修改为下面这样的，为了实现顶级菜单
// see: https://github.com/arco-design/arco-design-pro-vue/issues/294
const appClientMenus = mixinRoutes.map((el) => {
  let { name } = el;
  const { path, meta, redirect, children } = el;
  if (meta?.hideChildrenInMenu && children?.length) {
    name = children[0].name;
  }
  return {
    name,
    path,
    meta,
    redirect,
    children,
  };
});

export default appClientMenus;
