import {Link as RouterLink, useNavigate} from 'react-router-dom';
import {Box, Link, Container, Typography} from '@mui/material';
import AuthLayout from '../../layouts/AuthLayout';
import RegisterForm from './RegisterForm';
import {RootStyle, SectionStyle, ContentStyle} from "../RootStyle";
import {setLocalStorage} from "../../helpers/localStorage";

export default function Register(props) {
    const navigate = useNavigate();

    const handleSubmitSuccess = (data) => {
        setLocalStorage('session_id', data.token);
        props.setUser(data.user);
        navigate('/', {replace: true});
    };

    return (
        <RootStyle title="Register | def-project">
            <AuthLayout>
                Already have an account? &nbsp;
                <Link underline="none" variant="subtitle2" component={RouterLink} to="/login">
                    Login
                </Link>
            </AuthLayout>

            <SectionStyle sx={{display: {xs: 'none', md: 'flex'}}}>
                <Typography variant="h3" sx={{px: 5, mt: 10, mb: 5}}>
                    Manage the job more effectively with Minimal
                </Typography>
                <img alt="register" src="/static/illustrations/illustration_register.png"/>
            </SectionStyle>

            <Container>
                <ContentStyle>
                    <Box sx={{mb: 5}}>
                        <Typography variant="h4" gutterBottom>
                            Get started absolutely free.
                        </Typography>
                    </Box>

                    <RegisterForm onSubmitSuccess={handleSubmitSuccess}/>

                    <Typography
                        variant="subtitle2"
                        sx={{
                            mt: 3,
                            textAlign: 'center',
                            display: {sm: 'none'}
                        }}
                    >
                        Already have an account?&nbsp;
                        <Link underline="hover" to="/login" component={RouterLink}>
                            Login
                        </Link>
                    </Typography>
                </ContentStyle>
            </Container>
        </RootStyle>
    );
}
