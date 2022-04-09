import * as Yup from 'yup';
import {useState} from 'react';
import {Link as RouterLink} from 'react-router-dom';
import {useFormik, Form, FormikProvider} from 'formik';

import {
    Link,
    Stack,
    Checkbox,
    TextField,
    IconButton,
    InputAdornment,
    FormControlLabel
} from '@mui/material';
import {LoadingButton} from '@mui/lab';

import Iconify from '../../components/Iconify';
import axios from "../../plugins/axios";
import urls from "../../urls"

export default function LoginForm(props) {
    const {onLoginSuccess} = props;
    const [showPassword, setShowPassword] = useState(false);

    const LoginSchema = Yup.object().shape({
        username: Yup.string().required('Username is required'),
        password: Yup.string().required('Password is required')
    });

    const login = (form, setStatus, setErrors) => {
        axios("post", urls.auth, form).then((response => {
            let data = response.data

            if (data.error) {
                setStatus({success: false});
                setErrors({username: data["details"]});

                return null
            }
            onLoginSuccess(data)
        })).catch((error => {
            setStatus({success: false});
            if (error.response === undefined) {
                setErrors({username: "Network error"});
            } else {
                setErrors({username: error.response.data.message});
            }
        }))
    }

    const formik = useFormik({
        initialValues: {
            username: '',
            password: '',
            remember: true
        },
        validationSchema: LoginSchema,
        onSubmit: async (values, {
            setErrors,
            setStatus
        }) => {
            try {
                login(values, setStatus, setErrors)
            } catch (error) {
                const message = (error.response && error.response.data.message) || 'Something went wrong';

                setStatus({success: false});
                setErrors({submit: message});
            }
        }
    });

    const {errors, touched, values, isSubmitting, handleSubmit, getFieldProps} = formik;

    const handleShowPassword = () => {
        setShowPassword((show) => !show);
    };

    return (
        <FormikProvider value={formik}>
            <Form autoComplete="off" noValidate onSubmit={handleSubmit}>
                <Stack spacing={3}>
                    <TextField
                        fullWidth
                        autoComplete="username"
                        type="text"
                        label="Username"
                        {...getFieldProps('username')}
                        error={Boolean(touched.username && errors.username)}
                        helperText={touched.username && errors.username}
                    />

                    <TextField
                        fullWidth
                        autoComplete="current-password"
                        type={showPassword ? 'text' : 'password'}
                        label="Password"
                        {...getFieldProps('password')}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton onClick={handleShowPassword} edge="end">
                                        <Iconify icon={showPassword ? 'eva:eye-fill' : 'eva:eye-off-fill'}/>
                                    </IconButton>
                                </InputAdornment>
                            )
                        }}
                        error={Boolean(touched.password && errors.password)}
                        helperText={touched.password && errors.password}
                    />
                </Stack>

                <Stack direction="row" alignItems="center" justifyContent="space-between" sx={{my: 2}}>
                    <FormControlLabel
                        control={<Checkbox {...getFieldProps('remember')} checked={values.remember}/>}
                        label="Remember me"
                    />

                    <Link component={RouterLink} variant="subtitle2" to="#" underline="hover">
                        Forgot password?
                    </Link>
                </Stack>

                <LoadingButton
                    fullWidth
                    size="large"
                    type="submit"
                    variant="contained"
                    loading={isSubmitting}
                >
                    Login
                </LoadingButton>
            </Form>
        </FormikProvider>
    );
}
