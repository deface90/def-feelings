import React from 'react';
import {Outlet} from 'react-router-dom';
// material
import {styled} from '@mui/material/styles';
//
import {compose} from 'recompose';
import {connect} from 'react-redux';
import {setUser} from '../../components/App/AppState';

import DashboardNavbar from './DashboardNavbar';
import DashboardSidebar from './DashboardSidebar';

// ----------------------------------------------------------------------

const APP_BAR_MOBILE = 64;
const APP_BAR_DESKTOP = 92;

const RootStyle = styled('div')({
    display: 'flex',
    minHeight: '100%',
    overflow: 'hidden'
});

const MainStyle = styled('div')(({theme}) => ({
    flexGrow: 1,
    overflow: 'auto',
    minHeight: '100%',
    paddingTop: APP_BAR_MOBILE + 24,
    paddingBottom: theme.spacing(10),
    [theme.breakpoints.up('lg')]: {
        paddingTop: APP_BAR_DESKTOP + 24,
        paddingLeft: theme.spacing(2),
        paddingRight: theme.spacing(2)
    }
}));

// ----------------------------------------------------------------------

class AppLayout extends React.Component {
    constructor(props) {
        super(props);
        this.user = this.props.user;
    }

    render() {
        return (
            <RootStyle>
                <DashboardNavbar user={this.user}/>
                <DashboardSidebar user={this.user}/>
                <MainStyle>
                    <Outlet/>
                </MainStyle>
            </RootStyle>
        );
    }
}

export default compose(
    connect(
        state => ({
            app: state.app,
            user: state.app.user,
            login: state.login,
        }),
        dispatch => ({
            setUser: (user) => dispatch(setUser(user)),
        }),
    ),
)(AppLayout);