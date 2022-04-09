import React from "react";
import axios from "../../plugins/axios";
import urls from "../../urls";
import Page from "../../components/Page";
import {
    Box, Button,
    Card,
    CardContent,
    CardHeader,
    Container,
    Divider,
    Grid,
    Stack,
    Typography
} from "@mui/material";
import PropTypes from "prop-types";
import {fToNow} from "../../utils/formatTime";
import StatusFilter from "./StatusFilters";

export default class StatusList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isLoading: true,
            statuses: [],
            filters: {
                page: 0,
                limit: 5,
                feelings: [],
                datetime_start: null,
                datetime_end: null
            },
            hasMore: false
        }
    }

    componentDidMount() {
        this.loadStatuses();
    }

    loadStatuses = (clean) => {
        if (clean) {
            this.setState({statuses: []});
        }
        axios("post", urls["status-list"], this.state.filters).then(response => {
            if (response.data === null) {
                this.setState({hasMore: false});
                return;
            }
            if (response.data.error !== undefined && response.data.error !== "") {
                console.log(response.data.error);
                return;
            }
            let statuses = this.state.statuses;
            statuses.push(...response.data);
            this.setState({statuses: statuses})

            const totalCount = parseInt(response.headers['X-Pagination-Total-Count']);
            if (totalCount > (this.state.filters.page + 1) * this.state.filters.limit) {
                this.setState({hasMore: true});
            }
        }).finally(() => {
            this.setState({isLoading: false})
        });
    }

    nextPage = () => {
        let filters = this.state.filters;
        filters.page++;
        this.setState({filters: filters});
        this.loadStatuses(false);
    }

    setFilterFeelings = (value) => {
        let filters = this.state.filters;
        filters.feelings = value;
        filters.page = 0;
        this.setState({filters: filters});
        this.loadStatuses(true);
    }
    setFilterStartDate = (value) => {
        let filters = this.state.filters;
        filters.datetime_start = value.toISOString();
        filters.page = 0;
        this.setState({filters: filters});
        this.loadStatuses(true);
    }
    setFilterEndDate = (value) => {
        let filters = this.state.filters;
        filters.datetime_end = value.toISOString();
        filters.page = 0;
        this.setState({filters: filters});
        this.loadStatuses(true);
    }

    render() {
        return (
            <Page title="Your statuses | def-feelings">
                <Container>
                    <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
                        <Typography variant="h4">Your statuses</Typography>
                    </Stack>


                    <StatusFilter
                        filterStartDate={this.state.filters.datetime_start}
                        filterEndDate={this.state.filters.datetime_end}
                        onChangeFeelingsFilter={(newValue) => {
                            this.setFilterFeelings(newValue);
                        }}
                        onChangeStartDateFilter={(newValue) => {
                            this.setFilterStartDate(newValue);
                        }}
                        onChangeEndDateFilter={(newValue) => {
                            this.setFilterEndDate(newValue);
                        }}
                    />

                    <Grid container spacing={3}>
                        {this.state.isLoading ? "" : this.state.statuses.map((status) => (
                            <Grid item xs={12} key={status.id}>
                                <StatusCard key={status.id} status={status}/>
                            </Grid>
                        ))}

                        <Divider/>
                        {!this.state.hasMore ? "" :
                            <Grid item xs={12} key={0}>
                                <Button fullWidth size="large" color="warning" onClick={this.nextPage}>
                                    Load more
                                </Button>
                            </Grid>
                        }
                    </Grid>
                </Container>
            </Page>
        )
    }
}

StatusCard.propTypes = {
    status: PropTypes.shape({
        message: PropTypes.string,
        created: PropTypes.string,
        feelings: PropTypes.array,
    }),
};

function StatusCard({status}) {
    const {feelings, message, created} = status;
    const parsedCreated = Date.parse(created);
    const parsedFeelings = feelings.join(", ");

    return (
        <Card>
            <CardHeader title={parsedFeelings} sx={{mb: 3}}/>
            <CardContent>
                {message}
            </CardContent>
            <Divider/>

            <Box sx={{p: 2, textAlign: 'right'}}>
                {fToNow(parsedCreated)}
            </Box>
        </Card>
    );
}
