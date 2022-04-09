import {useNavigate} from "react-router";

export const withRouter = (Component) => {
    return (props) => {
        const history = useNavigate();
        return <Component history={history} {...props} />;
    };
};