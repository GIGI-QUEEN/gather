import { Routes, Route } from 'react-router-dom';
import { PostPage } from '../../pages/post/Post';
import { UserPage } from '../../pages/user/User';
import CreatePostPage from '../../pages/create/Create';
import Me from '../../pages/me/Me';
import Users from '../../pages/users/Users';
import { Feed } from '../../pages/posts/Posts';
import Followers from '../../pages/followers/Followers';
import Followings from '../../pages/followings/Followings';
import SettingsPage from '../../pages/settings/Settings';
import GroupPage from '../../pages/group/GroupPage';
import GroupsPage from '../../pages/groups/Groups';
import ChatPage2 from '../../pages/ws-chat/ChatPage';
import NotificationsCenter from '../../pages/notificationsCenter/NotificationsCenter';
import GroupSettingsPage from '../../pages/groupSettings/GroupSettings';

const AllRoutes = () => {
  return (
    <div>
      <Routes>
        <Route path="chat" element={<ChatPage2 />} />
        <Route path="posts" element={<Feed />} />
        <Route path="/" element={<Feed />} />
        <Route path="post/:id" element={<PostPage fetch_url={'post'} />} />
        <Route path="create" element={<CreatePostPage />} />
        <Route path="user/:id" element={<UserPage />} />
        <Route path="user/all" element={<Users />} />
        <Route path="me" element={<Me />} />
        <Route path="me/settings" element={<SettingsPage />} />
        <Route path="user/:id/followers" element={<Followers />} />
        <Route path="user/:id/followings" element={<Followings />} />
        <Route path="group/:id" element={<GroupPage />} />
        <Route path="groups" element={<GroupsPage />} />
        <Route
          path="/group/:groupid/post/:id"
          element={<PostPage fetch_url={'/group/post'} />}
        />
        <Route
          path="/group/:groupid/settings"
          element={<GroupSettingsPage />}
        />
        <Route path="notifications" element={<NotificationsCenter />} />
      </Routes>
    </div>
  );
};

export default AllRoutes;
