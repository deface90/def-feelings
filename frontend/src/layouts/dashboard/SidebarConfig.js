// component
import Iconify from '../../components/Iconify';

// ----------------------------------------------------------------------

const getIcon = (name) => <Iconify icon={name} width={22} height={22} />;

const sidebarConfig = [
  {
    title: 'create status',
    path: '/status/create',
    icon: getIcon('eva:activity-fill')
  },
  {
    title: 'status list',
    path: '/status/list',
    icon: getIcon('eva:file-text-fill')
  },
  {
    title: 'profile edit',
    path: '/profile/edit',
    icon: getIcon('eva:edit-2-fill')
  }
];

export default sidebarConfig;
