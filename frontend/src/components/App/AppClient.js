import {Provider} from 'react-redux';
import {ReactNotifications} from 'react-notifications-component';
import Navigation from '../../routes';
import {BaseOptionChartStyle} from '../charts/BaseOptionChart';
import {Component} from "react";

import {createBrowserHistory} from 'history';
import storeRedux from '../../redux/store';
import ThemeProvider from "../../theme";
import ThemeColorPresets from "../ThemeColorPresets";

export default class App extends Component {
    render() {
        const history = createBrowserHistory();
        return (
            <Provider store={storeRedux}>
                <ThemeProvider>
                    <ThemeColorPresets>
                        <BaseOptionChartStyle/>
                        <ReactNotifications/>
                        <Navigation history={history}/>
                    </ThemeColorPresets>
                </ThemeProvider>
            </Provider>
        );
    }
}
