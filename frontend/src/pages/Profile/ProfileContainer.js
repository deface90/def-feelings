import { compose } from 'recompose';
import { connect } from 'react-redux';

import Profile from './Profile';

import {
    setUser
} from '../../components/App/AppState';

export default compose(
    connect(
        state => ({
            user: state.app.user
        }),
        dispatch => ({
            setUser: (user) => dispatch(setUser(user))
        }),
    ),
)(Profile);
