import React from 'react';
import {Store} from "react-notifications-component"
import {Container, Stack, Typography} from "@mui/material";
import Page from "../../components/Page";
import ProfileForm from "./ProfileForm";

export default function Profile(props) {
    const {user, setUser} = props;

    const handleSubmitSuccess = (user) => {
        setUser(user);
        Store.addNotification({
            title: 'Success',
            message: 'Profile saves',
            type: 'success',
            insert: 'top',
            container: 'top-right',
            dismiss: {
                duration: 3000,
                onScreen: false,
                pauseOnHover: true
            }
        });
    };

    return (
        <Page title="Profile edit | def-feelings">
            <Container>
                <Stack direction="row" alignItems="center" justifyContent="space-between" mb={5}>
                    <Typography variant="h4">Profile edit</Typography>
                </Stack>
                <ProfileForm initValues={user} onSubmitSuccess={handleSubmitSuccess}/>
            </Container>
        </Page>
    );
}