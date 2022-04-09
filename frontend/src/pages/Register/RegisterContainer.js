import { compose } from 'recompose';
import { connect } from 'react-redux';

import Register from './Register';

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
)(Register);
