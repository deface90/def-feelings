import * as Yup from 'yup';
import {useState} from 'react';
import {Store} from 'react-notifications-component';
import {useFormik, Form, FormikProvider} from 'formik';

import {
    Stack,
    TextField,
    IconButton,
    InputAdornment,
    CardHeader,
    CardContent,
    RadioGroup,
    FormControlLabel, Radio, Box, Typography, Card
} from '@mui/material';
import {LoadingButton} from '@mui/lab';

import axios from "../../plugins/axios";
import urls from "../../urls"
import Iconify from '../../components/Iconify';

export default function RegisterForm(props) {
    const {onSubmitSuccess} = props;
    const [showPassword, setShowPassword] = useState(false);

    const saveProfile = (values, setStatus, setErrors, setSubmitting) => {
        values["notification_type"] = parseInt(values["notification_type"]);
        axios("post", urls["user-create"], values).then((response => {
            setSubmitting(false);
            let data = response.data

            if (data.error) {
                setStatus({success: false});
                Store.addNotification({
                    title: 'Registration error',
                    message: data.error,
                    type: 'danger',
                    insert: 'top',
                    container: 'top-right',
                    dismiss: {
                        duration: 3000,
                        onScreen: false,
                        pauseOnHover: true
                    }
                });

                return null
            }
            onSubmitSuccess(data);
        })).catch((error) => {
            console.log(error);
        })
    }

    const RegisterSchema = Yup.object().shape({
        username: Yup.string()
            .min(3, 'Too Short!')
            .max(50, 'Too Long!')
            .required('Username required'),
        full_name: Yup.string().min(2, 'Too Short').max(50, 'Too Long'),
        email: Yup.string().email('Email must be a valid email address'),
        tg_username: Yup.string().min(2, 'Too Short').max(50, 'Too Long'),
        password: Yup.string().min(6, 'Too short').required('Password is required')
    });

    const formik = useFormik({
        initialValues: {
            username: '',
            full_name: '',
            email: '',
            password: '',
            tg_username: '',
            notification_type: 0
        },
        validationSchema: RegisterSchema,
        onSubmit: (values, {
            setErrors,
            setStatus,
            setSubmitting
        }) => {
            try {
                saveProfile(values, setStatus, setErrors, setSubmitting);
            } catch (error) {
                const message = (error.response && error.response.data.message) || 'Что-то пошло не так';

                setStatus({success: false});
                setErrors({submit: message});
            }
        }
    });

    const {errors, touched, handleSubmit, isSubmitting, getFieldProps} = formik;

    const handleShowPassword = () => {
        setShowPassword((show) => !show);
    };

    const NOTIFICATION_TYPES = [
        {
            value: 0,
            title: 'None',
            description: "Disable notifications"
        },
        {
            value: 1,
            title: 'Browser',
            description: "Notify you by browser native notifications"
        },
        {
            value: 2,
            title: 'Telegram',
            description: "Notify you by telegram bot"
        }
    ];

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
                        autoComplete="full-name"
                        type="text"
                        label="Full name"
                        {...getFieldProps('full_name')}
                        error={Boolean(touched.full_name && errors.full_name)}
                        helperText={touched.full_name && errors.full_name}
                    />

                    <TextField
                        fullWidth
                        autoComplete="email"
                        type="email"
                        label="Email address"
                        {...getFieldProps('email')}
                        error={Boolean(touched.email && errors.email)}
                        helperText={touched.email && errors.email}
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

                    <TextField
                        fullWidth
                        autoComplete="tg-username"
                        type="text"
                        label="Telegram username"
                        {...getFieldProps('tg_username')}
                        error={Boolean(touched.tg_username && errors.tg_username)}
                        helperText={touched.tg_username && errors.tg_username}
                    />

                    <Card sx={{my: 3}}>
                        <CardHeader title="Notification options"/>
                        <CardContent>
                            <RadioGroup
                                row
                                {...getFieldProps('notification_type')}
                            >
                                {NOTIFICATION_TYPES.map((option) => {
                                    const { value, title, description } = option;
                                    return (
                                        <FormControlLabel
                                            key={value}
                                            value={value}
                                            control={<Radio />}
                                            label={
                                                <Box sx={{ml: 1}}>
                                                    <Typography variant="subtitle2">{title}</Typography>
                                                    <Typography variant="body2" sx={{color: 'text.secondary'}}>
                                                        {description}
                                                    </Typography>
                                                </Box>
                                            }
                                            sx={{flexGrow: 1, py: 3}}
                                        />
                                    )})}
                            </RadioGroup>
                        </CardContent>
                    </Card>

                    <LoadingButton
                        fullWidth
                        size="large"
                        type="submit"
                        variant="contained"
                        loading={isSubmitting}
                    >
                        Register
                    </LoadingButton>
                </Stack>
            </Form>
        </FormikProvider>
    );
}
