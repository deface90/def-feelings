import {clearLocalStorage} from '../../helpers/localStorage';
import {useNavigate} from 'react-router-dom';

function Logout (props) {
    const navigate = useNavigate();

    props.setUser(null);
    clearLocalStorage('session_id');
    navigate('/', { replace: true });

    return null;
}

export default Logout;
