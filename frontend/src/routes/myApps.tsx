import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import {styled} from "@mui/material/styles";
import AddRoundedIcon from '@mui/icons-material/AddRounded';
import SyncAltRoundedIcon from '@mui/icons-material/SyncAltRounded';
import {AppPanel, AppBundle} from "./appPanel.tsx";
import {useLoaderData, useNavigate, useParams} from "react-router-dom";
import React, {Dispatch, MutableRefObject, SetStateAction, useRef, useState} from "react";
import Alert from "@mui/material/Alert";
import {Snackbar} from "@mui/material";
import {InstallationConfigPopUp, InstallationProgressPopup} from "./installDetail.tsx";
import {fetchEventSource} from "@microsoft/fetch-event-source";

export const AppIcon = styled('div')({
    width: 70,
    height: 70,
    objectFit: 'cover',
    overflow: 'hidden',
    flexShrink: 0,
    borderRadius: 20,
    alignItems:'center',
    backgroundColor: 'rgba(0,0,0,0.08)',
    '& > img': {
        width: '100%',
        display: "block",
        margin:"auto",
        padding:"auto",
    },
});

export async function AppListLoader({params}: { params: { udid: string }} )  {
    return fetch(`/api/devices/${params.udid}`)
        .then(function (response) {
            if (!response.ok) {
                throw Error("SERVERR");
            }
            return response.json();
        }).then(function (nodes: AppBundle[]) {
        return nodes
    }).catch(function (error) {
        throw Error(error.toString());
    })
}

export function MyApps() {
    const { udid } = useParams();
    const apps  = useLoaderData() as AppBundle[];
    const navigate = useNavigate()
    const inputFile = useRef<HTMLInputElement | null>(null);
    const [openAlert, setOpenAlert] = useState(false);
    const [alertMsg, setAlertMsg] = useState("");

    const [openProgressPopUps, setOpenProgressPopUps] = useState(false);
    const [displayLog, setDisplayLog] = useState(["installation starts."]);
    const [finished, setFinished] = useState(false);
    const [success, setSuccess] = useState(false);
    const [appName, setAppName] = React.useState("");

    const [openConfigPopUps, setOpenConfigPopUp] = useState(false);
    const [ipaHash, setIpaHash] = React.useState("");



    if(udid == undefined){
        throw Error("undefined udid")
    }
    const nodeList = apps
        .map((node: AppBundle) => (
            <AppPanel app={node} key={node.BundleIdentifier} udid={udid}/>
        ))

    const onButtonClick = () => {
        inputFile?.current?.click();
    };
    const handleClose = () => {
        setOpenAlert(false);
    };

    return (
        <div style={{
            padding: "10px",
            maxWidth: "430px",
            marginLeft: "auto",
            marginRight: "auto",
            width: "100%",
            height: "100%",
            boxSizing: "border-box"
        }}>

            <div style={{display: "flex", justifyContent: "space-between"}}>
                <SyncAltRoundedIcon sx={{color: "#0f7e82"}} fontSize="large" style={{cursor: "pointer"}}
                                    onClick={() => navigate("/devices")}/>
                <input type='file' id='file' ref={inputFile} style={{display: 'none'}} onChange={async (file) => {
                    const packageInfo = await uploadIpa(file, setOpenAlert, setAlertMsg, inputFile)
                    if(packageInfo != null) {
                        if (!packageInfo.contains_plug_in) {
                            doInstall(packageInfo.md5, false, udid, packageInfo.cf_bundle_name, setOpenProgressPopUps, setDisplayLog, setFinished, setSuccess, setAppName)
                        }else {
                            setAppName(packageInfo.cf_bundle_name)
                            setIpaHash(packageInfo.md5)
                            setOpenConfigPopUp(true)
                        }
                    }
                }}/>
                <AddRoundedIcon sx={{color: "#0f7e82"}} style={{cursor: "pointer"}} fontSize="large" onClick={onButtonClick}/>
            </div>
            <Typography fontSize={'2.5rem'} fontFamily={"CircularSpBold"}>My Apps</Typography>
            <div style={{display: "flex", flexDirection: "row", justifyContent: "space-between", alignItems: "center"}}>
                <Typography fontSize={'1.5rem'} fontFamily={"CircularSpBold"}>Installed</Typography>
                    <Typography color={"#0f7e82"} fontSize={'1rem'} fontFamily={"CircularSpBold"} style={{cursor:"pointer"}}>Refresh All </Typography>
                </div>
                <Grid container style={{paddingBottom: "10px"}} direction={"column"}>
                    {nodeList}
                </Grid>
            <Snackbar open={openAlert} autoHideDuration={5000} onClose={handleClose}>
                <Alert variant="filled" severity="error">
                    {alertMsg}
                </Alert>
            </Snackbar>
            <InstallationProgressPopup app={appName} openPopUps={openProgressPopUps} setOpenPopUps={setOpenProgressPopUps} key={"install-progress"} log={displayLog} finished={finished} success={success}/>
            <InstallationConfigPopUp app={appName} openPopUps={openConfigPopUps} setOpenPopUps={setOpenConfigPopUp} key={"install-config"} onInstall={(removePlugIns:boolean)=>{
                doInstall(ipaHash, removePlugIns, udid, appName, setOpenProgressPopUps, setDisplayLog, setFinished, setSuccess, setAppName)
            }}/>

        </div>
    )
}

async function uploadIpa(
    file:React.ChangeEvent<HTMLInputElement>,
    setOpenAlert: Dispatch<SetStateAction<boolean>>,
    setAlertMsg: Dispatch<SetStateAction<string>>,
    inputFile:MutableRefObject<HTMLInputElement|null>
):Promise<packageUploadResponse|null>{
    const  files = file.target?.files
    if(files) {
        const binary = files[0]
        const formData = new FormData();
        formData.append("file", binary);
        const packageInfo = await fetch("/api/packages", {
            method: "POST",
            body: formData,
        }).then(function (response) {
            return response.json().then(data => ({ status: response.status, body: data }))
        }).then(function (data:{status:number, body:packageUploadResponse}) {
            if(data.status != 200){
                throw Error(data.body.error_msg)
            }
            return data.body
        }).catch(function (error){
            setOpenAlert(true)
            setAlertMsg(error.toString())
            return null
        })
        if (inputFile.current) {
            inputFile.current.value = "";
            inputFile.current.type = "file";
        }

        if(packageInfo != null){
            return packageInfo
        }
    }
    return null
}

type packageUploadResponse = {
    cf_bundle_name: string
    cf_bundle_short_version_string: string
    cf_bundle_identifier: string
    contains_plug_in: boolean
    md5: string
    error_msg:string
}

let previous = ""
export function doInstall(
    ipa: string,
    removePlugIns: boolean,
    udid: string,
    appName:string,
    setOpenPopUps: React.Dispatch<React.SetStateAction<boolean>>,
    setDisplayLog: React.Dispatch<React.SetStateAction<string[]>>,
    setFinished: React.Dispatch<React.SetStateAction<boolean>>,
    setSuccess: React.Dispatch<React.SetStateAction<boolean>>,
    setAppName: React.Dispatch<React.SetStateAction<string>>

){
    {
        setAppName(appName)
        setOpenPopUps(true)
        setDisplayLog(["installation starts."])
        setFinished(false)
        previous = ""
        const data = {
            "package": ipa,
            "remove_plug_ins": removePlugIns
        }
        fetchEventSource(`/api/devices/${udid}`, {
            method: 'POST',
            headers: {
                "Content-Type": 'application/json',
            },
            body: JSON.stringify(data),
            onmessage(event) {
                if(event.data == "ok"){
                    setFinished(true)
                    setSuccess(true)
                }if(event.data == "failed"){
                    setFinished(true)
                    setSuccess(false)
                }else{
                    if(event.data != previous){
                        setDisplayLog(displayLog => [...(displayLog.slice(-999)), event.data])
                        previous = event.data
                    }

                }
            },
            onerror() {
                setFinished(true)
                setSuccess(false)
            }
        })
    }
}