import {Box, Grid, Typography} from "@mui/material";
import {RootStyle} from "./RootStyle";

export default function Loading() {
    return (
        <RootStyle title="Loading | def-project">
            <Grid
                container
                spacing={0}
                align="center"
                justify="center"
                direction="column"
            >
                <Grid item>
                    <Box sx={{maxWidth: 480, margin: 'auto', textAlign: 'center'}}>
                        <Box component="img" src="/static/logo.png"/>
                        <Typography variant="h3" paragraph>
                            Loading...
                        </Typography>
                        <Typography sx={{color: 'text.secondary'}}>
                            Application is loading, please wait
                        </Typography>
                    </Box>
                </Grid>
            </Grid>
        </RootStyle>
    )
}