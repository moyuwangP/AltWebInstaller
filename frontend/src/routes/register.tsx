import Typography from "@mui/material/Typography";
import InstallMobileIcon from '@mui/icons-material/InstallMobile';
import {Button} from "@mui/material";
import {useNavigate} from "react-router-dom";


export function Register(){
    const navigate = useNavigate()

    return (
        <div style={{
            display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                flexDirection:"column",
                width: "100vw",
                height: "90dvh"
            }}>
            <div style={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                flexDirection:"column",
                width: "100vw",
                height: "80vh"
            }}>
                <Typography
                    fontSize={'2.5rem'}
                    sx={{color: "#0f7e82"}}
                    fontFamily={"CircularSpBold"}>
                    AltWebInstaller
                </Typography>
                <Typography fontFamily={"CircularSpBold"}>Download Profile to Register Your Device</Typography>
                <Button href="/api/udid/registration-file">
            <InstallMobileIcon
                sx={{color: "#0f7e82"}}
                style={{padding:"16", cursor:"pointer"}}
                fontSize="large"/>
                </Button>
            </div>
            <Typography
                sx={{color: "#0f7e82"}}
                fontFamily={"CircularSpBold"}
                fontSize={'1.25rem'}
                onClick={() => navigate("/devices")}
                style={{cursor:"pointer"}}
            >
                paired devices
            </Typography>
        </div>
    )
}