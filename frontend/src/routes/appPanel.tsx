import {useEffect, useState} from "react";
import {FinalColor} from "extract-colors/lib/types/Color";
import {extractColors} from "extract-colors";
import Grid from "@mui/material/Grid";
import {Avatar} from "@mui/material";
import BrokenImageRoundedIcon from "@mui/icons-material/BrokenImageRounded";
import Typography from "@mui/material/Typography";
import {AppIcon} from "./myApps.tsx";
import {styled} from "@mui/material/styles";
import {InstallationProgressPopup} from "./installDetail.tsx";
import {fetchEventSource} from "@microsoft/fetch-event-source";

type Props = {
    app:AppBundle
    udid:string
    key:string
}
export type AppBundle ={
    BundleIdentifier: string
    DisplayName: string
    RefreshedAt: string
    Version: string
    ipa_hash:string
    RemovePlugIns:boolean
}
export function AppPanel(prop: Props){
    const src = `/api/packages/${prop.app.ipa_hash}/app-icon`
    const [accentColor, setAccentColor] = useState<FinalColor>({
        area: 0, hex: "", hue: 0, intensity: 0, lightness: 0, saturation: 0,
        red:1, blue:1, green:1});
    const [appIcon, setAppIcon] = useState("");

    useEffect(() => { extractColors(src, {crossOrigin:"anonymous"})
        .then(function (colors: FinalColor[]) {

            const c = colors[0]
            if(c!= null){
                setAccentColor(c)
            }
            setAppIcon(src)
        })
        .catch(console.error)}, []);

    let displayIdentifier = prop.app.BundleIdentifier
    if (displayIdentifier.length >26){
        displayIdentifier = displayIdentifier.substring(0,23)+"..."
    }
    return (
        <Grid item>
            <Item c={accentColor}>
                <div style={{display:"flex", flexDirection:"row", alignItems:"center"}}>
                    <AppIcon>
                        <Avatar
                            variant="square"
                            sx={{ width: 70, height: 70}}
                            src={appIcon}
                            alt={prop.app.BundleIdentifier}
                        > <BrokenImageRoundedIcon style={{transform:"scale(5)"}}/></Avatar>
                    </AppIcon>
                    <div style={{paddingLeft:10}}>
                        <Typography fontSize={'1.25rem'} fontFamily={"CircularSpBold"}>{prop.app.DisplayName}</Typography>
                        <Typography fontSize={'0.75rem'} fontFamily={"CircularSp"} >{displayIdentifier}</Typography>
                    </div>
                </div>
                <RefreshBox app={prop.app} key={prop.app.BundleIdentifier} udid={prop.udid}></RefreshBox>
            </Item>
        </Grid>
    )
}

type ItemProp ={
    c:FinalColor
}
const Item = styled('div')((item:ItemProp) => ({
    borderRadius: 30,
    height: "7rem",
    backgroundColor: "rgba("+item.c.red+", "+item.c.green+", "+item.c.blue+", 0.25)",
    marginTop:16,
    display: 'flex',
    alignItems:'center',
    paddingLeft:16,
    paddingRight:16,
    justifyContent:'space-between'
}));

let previous = ""
function RefreshBox(prop: Props){
    const [openPopUps, setOpenPopUps] = useState(false);
    const [displayLog, setDisplayLog] = useState(["installation starts."]);
    const [finished, setFinished] = useState(false);
    const [success, setSuccess] = useState(false);


    const refreshedAt = Date.parse(prop.app.RefreshedAt) / 1000
    const difference = (Date.now() / 1000)-refreshedAt
    const daysRemaining = 7- Math.floor(difference/60/60/24)
    let expires = "EXPIRED"
    if(daysRemaining ==1){
        expires = daysRemaining +" DAY"
    }else if(daysRemaining > 1){
        expires = daysRemaining +" DAYS"
    }
    return (<div style={{display: "flex", flexDirection: "column", alignItems: "center"}}>
        <Typography fontSize={'0.75rem'} fontFamily={"CircularSp"}>Expires in</Typography>
        <div style={{
            backgroundColor: "#13c204",
            borderRadius: 30,
            height: "2rem",
            width: "5rem",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            cursor: "pointer"
        }}
             onClick={() => {
                 setOpenPopUps(true)
                 setDisplayLog(["installation starts."])
                 setFinished(false)
                 // setProgress(-1)
                 previous = ""
                 const data = {
                     "package": prop.app.ipa_hash,
                     "remove_plug_ins": prop.app.RemovePlugIns
                 }
                 fetchEventSource(`/api/devices/${prop.udid}`, {
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
             }}>
            <Typography color="white" fontSize={'1rem'}
                        fontFamily={"CircularSpBold"}>{expires}</Typography>
        </div>
        <InstallationProgressPopup app={prop.app.DisplayName} openPopUps={openPopUps} setOpenPopUps={setOpenPopUps} key={prop.app.BundleIdentifier} log={displayLog} finished={finished} success={success}/>
    </div>)
}

