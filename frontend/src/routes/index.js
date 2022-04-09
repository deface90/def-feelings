import React, {useState, useEffect} from "react";
import {compose} from 'recompose';
import {connect} from 'react-redux';
import {getLocalStorage} from "../helpers/localStorage";
import {setUser} from "../components/App/AppState";
import {withRouter} from "../plugins/router";

import Main from "./main";
import Auth from "./auth";
import Loading from "../pages/Loading"
import axios from "../plugins/axios";

const Routes = (props) => {
    const {app, history} = props;
    const {user} = app;

    const [isLoading, setLoading] = useState(true);

    const handleChangeLocalStorage = () => {
        console.log('handleChangeLocalStorage')
    }

    const handleSetAccount = async () => {
        const sessionId = getLocalStorage('session_id');

        if (!sessionId) {
            setLoading(false);
            return false
        }

        const user = await axios('post', '/auth/session', {session_id: sessionId}).then((response) => {
            const data = response.data;
            if (!data.status) {
                return null
            }

            return data.user;
        }).catch((error) => {
            console.log(error);
        });

        props.setUser(user);
        setLoading(false);
    }

    useEffect(() => {
        (async () => {
            await handleSetAccount();
        })();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);
    useEffect(() => {
        window.addEventListener('user_clear', handleChangeLocalStorage);
    }, []);
    useEffect(() => history.listen(() => {
        const isAccount = Boolean(localStorage.getItem('session_id'));
        if (isAccount) {
            return null
        }
        props.setUser({});
    }), [history, props])

    if (isLoading) {
        return (
            <Loading/>
        )
    }

    let Navigation = Main;
    if (!user || !user?.id) {
        Navigation = Auth;
    }

    return (
        <Navigation/>
    )
}

const RoutesRouter = withRouter(Routes);

export default compose(
    connect(
        state => ({
            app: state.app
        }),
        dispatch => ({
            setUser: (user) => dispatch(setUser(user))
        }),
    ),
)(RoutesRouter);
