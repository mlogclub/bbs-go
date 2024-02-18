export interface UserState {
  id: number;
  nickname?: string;
  username?: string;
  avatar?: string;
  smallAvatar?: string;
  gender?: string;
  email?: string;
  description?: string;
  status: number;
  roles?: string[];
}
