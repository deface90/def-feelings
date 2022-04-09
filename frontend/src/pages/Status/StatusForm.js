import axios from "../../plugins/axios";
import urls from "../../urls";
import {Store} from "react-notifications-component";
import * as Yup from "yup";
import {useForm, Controller} from 'react-hook-form';
import {yupResolver} from '@hookform/resolvers/yup';
import {useEffect, useMemo, useState} from "react";
import {Chip, TextField, Autocomplete, Stack} from '@mui/material';
import {LoadingButton} from '@mui/lab';
import {
    FormProvider, RHFTextField,
} from '../../components/hook-form';

export default function StatusForm(props) {
    const [inputValue, setInputValue] = useState('');
    const [feelings, setFeelings] = useState([]);

    useEffect(() => {
        let active = true;
        (async () => {
            axios("post", urls["feelings-list"], {title: inputValue, limit: 10}).then((response => {
                let data = response.data;
                if (data === null) {
                    setFeelings([]);
                    return
                }
                if (data.error !== undefined) {
                    console.log(data.error)
                    return
                }
                let f = [];
                for (let i = 0; i < data.length; i++) {
                    f.push(data[i].title);
                }
                if (active) {
                    setFeelings(f);
                }
            })).catch((error) => {
                console.log(error);
            })
        })();
        return () => {
            active = false;
        }
    }, [inputValue]);

    const {onSubmitSuccess} = props;

    const saveStatus = (values) => {
        axios("post", urls["status-create"], values).then((response => {
            let data = response.data

            if (data.error) {
                Store.addNotification({
                    title: 'Creating status error',
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
            onSubmitSuccess();
        })).catch((error) => {
            Store.addNotification({
                title: 'Creating status error',
                message: error,
                type: 'danger',
                insert: 'top',
                container: 'top-right',
                dismiss: {
                    duration: 3000,
                    onScreen: false,
                    pauseOnHover: true
                }
            });
        })
    }

    const NewStatusSchema = Yup.object().shape({
        feelings: Yup.array().min(1, 'Feelings is required'),
        message: Yup.string()
    });
    const defaultValues = useMemo(
        () => ({
            feelings: [],
            message: ""
        }),
        // eslint-disable-next-line react-hooks/exhaustive-deps
        []
    );
    const methods = useForm({
        resolver: yupResolver(NewStatusSchema),
        defaultValues,
    });
    const {
        reset,
        watch,
        control,
        handleSubmit,
        formState: {isSubmitting},
    } = methods;
    const values = watch();
    const onSubmit = async () => {
        saveStatus(values);
        reset();
    }

    return (
        <FormProvider methods={methods} onSubmit={handleSubmit(onSubmit)}>
            <Stack spacing={3}>
                <Controller
                    name="feelings"
                    control={control}
                    render={({field}) => (
                        <Autocomplete
                            {...field}
                            multiple
                            freeSolo
                            onInputChange={(e, value) => {
                                setInputValue(value);
                            }}
                            onChange={(event, newValue) => field.onChange(newValue)}
                            options={feelings}
                            filterOptions={(x) => x}
                            renderTags={(value, getTagProps) =>
                                value.map((option, index) => (
                                    <Chip {...getTagProps({index})} key={option} size="small" label={option}/>
                                ))
                            }
                            renderInput={(params) =>
                                <TextField label="Feelings" helperText="Separate feelings by pressing Enter" {...params} />
                            }
                        />
                    )}
                />

                <RHFTextField
                    multiline
                    name="message"
                    label="Message"
                />

                <LoadingButton type="submit" variant="contained" size="large" loading={isSubmitting}>
                    Save
                </LoadingButton>
            </Stack>
        </FormProvider>
    )
}