export const REDIRECT_ROUTE_NAME = 'Redirect';

export const DEFAULT_ROUTE_NAME = 'Dashboard';

export const NOT_FOUND_ROUTE_NAME = 'notFound';

export const WHITE_LIST = [
  // { name: NOT_FOUND_ROUTE_NAME, children: [] },
  { name: 'login', children: [] },
];

export const NOT_FOUND = {
  name: NOT_FOUND_ROUTE_NAME,
};

export const DEFAULT_ROUTE = {
  title: '仪表盘',
  name: DEFAULT_ROUTE_NAME,
  fullPath: '/dashboard',
};
