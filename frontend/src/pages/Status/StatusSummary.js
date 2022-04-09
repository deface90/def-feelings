import React from 'react';
import PropTypes from 'prop-types';
import {Link as RouterLink} from 'react-router-dom';

import {Box, Button, Card, CardHeader, Divider, Stack, Typography} from '@mui/material';

import {fToNow} from '../../utils/formatTime';

import Iconify from '../../components/Iconify';
import Scrollbar from '../../components/Scrollbar';

export default function StatusSummary(props) {
    const {statuses} = props;

    return (
        <Card>
            <CardHeader title="Last statuses"/>

            <Scrollbar>
                <Stack spacing={3} sx={{p: 3, pr: 0}}>
                    {statuses.map((status) => (
                        <StatusItem key={status.id} status={status}/>
                    ))}
                </Stack>
            </Scrollbar>

            <Divider/>

            <Box sx={{p: 2, textAlign: 'right'}}>
                <Button
                    to="/status/list"
                    size="small"
                    color="inherit"
                    component={RouterLink}
                    endIcon={<Iconify icon={'eva:arrow-ios-forward-fill'}/>}
                >
                    View all
                </Button>
            </Box>
        </Card>
    )
}

StatusItem.propTypes = {
    status: PropTypes.shape({
        message: PropTypes.string,
        created: PropTypes.string,
        feelings: PropTypes.array,
    }),
};

function StatusItem({status}) {
    const {feelings, message, created} = status;
    const parsedCreated = Date.parse(created);
    const parsedFeelings = feelings.join(", ");

    return (
        <Stack direction="row" alignItems="center" spacing={2}>
            <Box sx={{minWidth: 500}}>
                <Typography variant="subtitle2" noWrap>
                    {parsedFeelings}
                </Typography>
                <Typography variant="body2" sx={{color: 'text.secondary'}} noWrap>
                    {message}
                </Typography>
            </Box>
            <Typography variant="caption" sx={{pr: 3, flexShrink: 0, color: 'text.secondary'}}>
                {fToNow(parsedCreated)}
            </Typography>
        </Stack>
    );
}
