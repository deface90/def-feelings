import { compose } from 'recompose';
import { connect } from 'react-redux';

import Logout from './Logout';

import {
    setUser
} from '../../components/App/AppState';

export default compose(
    connect(
        state => ({
            login: state.login
        }),
        dispatch => ({
            setUser: (user) => dispatch(setUser(user))
        }),
    ),
)(Logout);
