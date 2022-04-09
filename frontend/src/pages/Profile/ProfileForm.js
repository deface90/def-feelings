import * as Yup from 'yup';
import {useState} from 'react';
import {useFormik, Form, FormikProvider} from 'formik';

import {
    Box,
    Typography,
    Stack,
    TextField,
    IconButton,
    InputAdornment,
    Radio,
    RadioGroup,
    FormControlLabel,
    CardContent, CardHeader, Card
} from '@mui/material';
import {LoadingButton} from '@mui/lab';

import axios from "../../plugins/axios";
import urls from "../../urls"
import Iconify from '../../components/Iconify';
import {Store} from "react-notifications-component";

export default function ProfileForm(props) {
    const {initValues, onSubmitSuccess} = props;
    const [showPassword, setShowPassword] = useState(false);

    const save = (values, setStatus, setErrors, setSubmitting) => {
        values["notification_type"] = parseInt(values["notification_type"]);
        values["notification_frequency"] = parseInt(values["notification_frequency"]);
        axios("post", urls["user-edit"] + initValues.id, values).then((response => {
            setSubmitting(false);
            let data = response.data

            if (data.error) {
                setStatus({success: false});
                Store.addNotification({
                    title: 'Profile edit error',
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
            onSubmitSuccess(values);
        })).catch((error) => {
            console.log(error);
        })
    }

    const ProfileSchema = Yup.object().shape({
        username: Yup.string()
            .min(3, 'Too Short!')
            .max(50, 'Too Long!')
            .required('Username required'),
        full_name: Yup.string().min(2, 'Too Short').max(50, 'Too Long'),
        email: Yup.string().email('Email must be a valid email address'),
        tg_username: Yup.string().min(2, 'Too Short').max(50, 'Too Long'),
        password: Yup.string().min(6, 'Too short'),
        notification_frequency: Yup.number().integer().min(1, 'Must be more than 0')
    });

    const formik = useFormik({
        initialValues: initValues,
        validationSchema: ProfileSchema,
        onSubmit: (values, {
            setErrors,
            setStatus,
            setSubmitting
        }) => {
            try {
                save(values, setStatus, setErrors, setSubmitting);
            } catch (error) {
                const message = (error.response && error.response.data.message) || 'Something went wrong';

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

                    <TextField
                        fullWidth
                        type="text"
                        label="Notifications frequency (in minutes)"
                        {...getFieldProps('notification_frequency')}
                        error={Boolean(touched.notification_frequency && errors.notification_frequency)}
                        helperText={touched.notification_frequency && errors.notification_frequency}
                    />

                    <LoadingButton
                        fullWidth
                        size="large"
                        type="submit"
                        variant="contained"
                        loading={isSubmitting}
                    >
                        Save
                    </LoadingButton>
                </Stack>
            </Form>
        </FormikProvider>
    );
}
