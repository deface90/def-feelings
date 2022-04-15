import React from "react";
import axios from "../../plugins/axios";
import urls from "../../urls";
import Page from "../../components/Page";
import {
    Box, Button,
    Card,
    CardHeader,
    Container,
    Divider,
    Grid,
    Stack,
    Typography
} from "@mui/material";
import PropTypes from "prop-types";
import FeelingFilter from "./FeelingFilter";

export default class FeelingList extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isLoading: true,
            feelings: [],
            filters: {
                page: 0,
                limit: 50,
                datetime_start: null,
                datetime_end: null
            },
            hasMore: false
        }
    }

    componentDidMount() {
        this.loadFeelings();
    }

    loadFeelings = (clean) => {
        if (clean) {
            this.setState({feelings: []});
        }
        axios("post", urls["feelings-frequency"], this.state.filters).then(response => {
            if (response.data === null) {
                this.setState({hasMore: false});
                return;
            }
            if (response.data.error !== undefined && response.data.error !== "") {
                console.log(response.data.error);
                return;
            }
            let feelings = this.state.feelings;
            feelings.push(...response.data);
            this.setState({feelings: feelings})

            /*const totalCount = parseInt(response.headers['X-Pagination-Total-Count']);
            if (totalCount > (this.state.filters.page + 1) * this.state.filters.limit) {
                this.setState({hasMore: true});
            }*/
        }).finally(() => {
            this.setState({isLoading: false})
        });
    }

    nextPage = () => {
        let filters = this.state.filters;
        filters.page++;
        this.setState({filters: filters});
        this.loadFeelings(false);
    }

    setFilterStartDate = (value) => {
        let filters = this.state.filters;
        filters.datetime_start = value.toISOString();
        filters.page = 0;
        this.setState({filters: filters});
        this.loadFeelings(true);
    }
    setFilterEndDate = (value) => {
        let filters = this.state.filters;
        filters.datetime_end = value.toISOString();
        filters.page = 0;
        this.setState({filters: filters});
        this.loadFeelings(true);
    }

    render() {
        return (
            <Page title="Your frequent feelings | def-feelings">
                <Container>
                    <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
                        <Typography variant="h4">Your frequent feelings</Typography>
                    </Stack>


                    <FeelingFilter
                        filterStartDate={this.state.filters.datetime_start}
                        filterEndDate={this.state.filters.datetime_end}
                        onChangeStartDateFilter={(newValue) => {
                            this.setFilterStartDate(newValue);
                        }}
                        onChangeEndDateFilter={(newValue) => {
                            this.setFilterEndDate(newValue);
                        }}
                    />

                    <Grid container spacing={3}>
                        {this.state.isLoading ? "" : this.state.feelings.map((feeling) => (
                            <Grid item xs={12} key={feeling.feeling.id}>
                                <FeelingCard key={feeling.feeling.id} item={feeling}/>
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

FeelingCard.propTypes = {
    status: PropTypes.shape({
        message: PropTypes.string,
        created: PropTypes.string,
        feelings: PropTypes.array,
    }),
};

function FeelingCard({item}) {
    return (
        <Card>
            <CardHeader title={item.feeling.title} sx={{mb: 3}}/>
            <Divider/>

            <Box sx={{p: 2, textAlign: 'right'}}>
                Usage: {item.frequency} time(s)
            </Box>
        </Card>
    );
}
