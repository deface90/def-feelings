import {Autocomplete, Chip, Stack, TextField} from "@mui/material";
import {useEffect, useState} from "react";
import axios from "../../plugins/axios";
import urls from "../../urls";
import DatePicker from '@mui/lab/DatePicker';
import {LocalizationProvider} from "@mui/lab";
import AdapterDateFns from '@mui/lab/AdapterDateFns';

export default function StatusFilter({
                                         filterStartDate,
                                         filterEndDate,
                                         onChangeFeelingsFilter,
                                         onChangeStartDateFilter,
                                         onChangeEndDateFilter
                                     }) {
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

    return (
        <Stack
            spacing={2}
            direction={{xs: 'column', sm: 'row'}}
            alignItems={{sm: 'left'}}
            sx={{mb: 2}}
        >
            <LocalizationProvider dateAdapter={AdapterDateFns}>
                <Autocomplete
                    multiple
                    freeSolo
                    onInputChange={(e, value) => {
                        console.log(value);
                        setInputValue(value);
                    }}
                    size={"small"}
                    style={{width: 300}}
                    options={feelings}
                    filterOptions={(x) => x}
                    onChange={(event, newValue) => onChangeFeelingsFilter(newValue)}
                    renderTags={(value, getTagProps) =>
                        value.map((option, index) => (
                            <Chip {...getTagProps({index})} key={option} size="small" label={option}/>
                        ))
                    }
                    renderInput={(params) => <TextField label="Feelings" {...params} />}
                />
                <DatePicker
                    label="Start date"
                    onChange={onChangeStartDateFilter}
                    value={filterStartDate}
                    renderInput={(params) => (
                        <TextField
                            size={"small"}
                            {...params}
                            fullWidth
                            sx={{
                                maxWidth: {md: 200},
                            }}
                        />
                    )}
                />

                <DatePicker
                    label="End date"
                    onChange={onChangeEndDateFilter}
                    value={filterEndDate}
                    renderInput={(params) => (
                        <TextField
                            {...params}
                            size={"small"}
                            fullWidth
                            sx={{
                                maxWidth: {md: 200},
                            }}
                        />
                    )}
                />
            </LocalizationProvider>
        </Stack>
    )
}