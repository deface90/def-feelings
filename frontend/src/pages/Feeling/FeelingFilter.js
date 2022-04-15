import DatePicker from '@mui/lab/DatePicker';
import {LocalizationProvider} from "@mui/lab";
import AdapterDateFns from '@mui/lab/AdapterDateFns';
import {Stack, TextField} from "@mui/material";

export default function FeelingFilter({
                                         filterStartDate,
                                         filterEndDate,
                                         onChangeStartDateFilter,
                                         onChangeEndDateFilter
                                     }) {

    return (
        <Stack
            spacing={2}
            direction={{xs: 'column', sm: 'row'}}
            alignItems={{sm: 'left'}}
            sx={{mb: 2}}
        >
            <LocalizationProvider dateAdapter={AdapterDateFns}>
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