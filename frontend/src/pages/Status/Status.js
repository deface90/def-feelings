import React from 'react';
import {Store} from "react-notifications-component"
import {Container, Grid, Stack, Typography} from "@mui/material";
import Page from "../../components/Page";
import StatusForm from "./StatusForm";
import StatusSummary from "./StatusSummary";
import axios from "../../plugins/axios";
import urls from "../../urls";

export default class Status extends React.Component {
    constructor(props) {
        super(props);
        this.state = {statuses: []}
    }

    componentDidMount() {
        this.loadStatuses();
    }

    handleSubmitSuccess = () => {
        Store.addNotification({
            title: 'Success',
            message: "Status created",
            type: 'success',
            insert: 'top',
            container: 'top-right',
            dismiss: {
                duration: 3000,
                onScreen: false,
                pauseOnHover: true
            }
        });
        this.loadStatuses();
    };

    loadStatuses = () => {
        axios("post", urls["status-list"], {limit: 5}).then(response => {
            if (response.data === null) {
                this.setState({statuses: []});
                return;
            }
            if (response.data.error !== undefined && response.data.error !== "") {
                Store.addNotification({
                    title: 'Intenal error',
                    message: 'Failed to get last status list',
                    type: 'danger',
                    insert: 'top',
                    container: 'top-right',
                    dismiss: {
                        duration: 3000,
                        onScreen: false,
                        pauseOnHover: true
                    }
                });
                this.setState({statuses: []})
            }
            this.setState({statuses: response.data})
        }).finally(() => {
            this.setState({isLoading: false})
        });
    }

    render() {
        const statuses = this.state.statuses;
        return (
            <Page title="Create status | def-feelings">
                <Container>
                    <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
                        <Typography variant="h4">Describe your current feelings</Typography>
                    </Stack>
                    <Grid container spacing={3}>
                        <Grid item xs={12} md={6} lg={8}>
                            <StatusForm onSubmitSuccess={this.handleSubmitSuccess}/>
                        </Grid>
                        <Grid item xs={12} md={6} lg={8}>
                            <StatusSummary statuses={statuses}/>
                        </Grid>
                    </Grid>
                </Container>
            </Page>
        )
    }
}