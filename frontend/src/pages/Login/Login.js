import {Link as RouterLink, useNavigate} from 'react-router-dom';
import {RootStyle, SectionStyle, ContentStyle} from "../RootStyle";
import {Stack, Link, Container, Typography} from '@mui/material';
import {setLocalStorage} from '../../helpers/localStorage';
import AuthLayout from '../../layouts/AuthLayout';
import LoginForm from './LoginForm';

export default function Login(props) {
    const navigate = useNavigate();

    const handleSubmitSuccess = (data) => {
        setLocalStorage('session_id', data.token);
        props.setUser(data.user);
        navigate('/', {replace: true});
    };

    return (
        <RootStyle title="Login | def-project">
            <AuthLayout>
                Don’t have an account? &nbsp;
                <Link underline="none" variant="subtitle2" component={RouterLink} to="/register">
                    Get started
                </Link>
            </AuthLayout>

            <SectionStyle sx={{display: {xs: 'none', md: 'flex'}}}>
                <Typography variant="h3" sx={{px: 5, mt: 10, mb: 5}}>
                    Hi, Welcome Back
                </Typography>
                <img src="/static/illustrations/illustration_login.png" alt="login"/>
            </SectionStyle>

            <Container maxWidth="sm">
                <ContentStyle>
                    <Stack sx={{mb: 5}}>
                        <Typography variant="h4" gutterBottom>
                            Sign in
                        </Typography>
                    </Stack>

                    <LoginForm onLoginSuccess={handleSubmitSuccess}/>

                    <Typography
                        variant="body2"
                        align="center"
                        sx={{
                            mt: 3,
                            display: {sm: 'none'}
                        }}
                    >
                        Don’t have an account?&nbsp;
                        <Link variant="subtitle2" component={RouterLink} to="register" underline="hover">
                            Get started
                        </Link>
                    </Typography>
                </ContentStyle>
            </Container>
        </RootStyle>
    );
}
