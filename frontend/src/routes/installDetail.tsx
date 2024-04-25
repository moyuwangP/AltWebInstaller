import React, {useEffect, useRef, useState} from "react";
import {
    Button,
    Card,
    CardContent, Checkbox, Collapse,
    Dialog, FormControlLabel,
    IconButton,
    IconButtonProps,
    LinearProgress
} from "@mui/material";
import Typography from "@mui/material/Typography";
import {styled} from "@mui/material/styles";
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {useNavigate} from "react-router-dom";

type ProgressProps = {
    app:string
    openPopUps: boolean
    setOpenPopUps: React.Dispatch<React.SetStateAction<boolean>>
    finished: boolean
    success: boolean
    key:string
    log:string[]
}

interface ExpandMoreProps extends IconButtonProps {
    expand: boolean;
}

const ExpandMore = styled((props: ExpandMoreProps) => {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const {expand,...other } = props;
    return <IconButton {...other} />;
})(({ theme, expand }) => ({
    transform: !expand ? 'rotate(0deg)' : 'rotate(180deg)',
    transition: theme.transitions.create('transform', {
        duration: theme.transitions.duration.shortest,
    }),
}));
export function InstallationProgressPopup(prop: ProgressProps){
    const navigate = useNavigate()
    const [expanded, setExpanded] = React.useState(false);
    const [progress, setProgress] = React.useState(-1);

    useEffect(() => {
        const extractedProgress = extractProgress(prop.log.slice(-1)[0])
        if(extractedProgress != -1){
            setProgress(extractedProgress)
        }
    }, [prop.log]);

    const handleExpandClick = () => {
        setExpanded(!expanded);
    };

    return (
        <Dialog fullScreen={false} open={prop.openPopUps}>
            <Card style={{width:"300px"}}>
                <CardContent>
                    <Typography variant="h5" component="div">
                        {prop.app}
                    </Typography>
                    <Typography>
                        {prop.finished? (prop.success?"Installed":"Failed"):"Installing..."}
                    </Typography>
                    <LinearProgress color={prop.finished?(prop.success?"success":"error"):"primary"} variant={prop.finished||progress != -1 ? "determinate":"indeterminate"} value={progress}/>
                    <div style={{display:"flex", flexDirection:"row-reverse"}}>
                    <ExpandMore
                        expand={expanded}
                        onClick={handleExpandClick}
                    >
                        <ExpandMoreIcon />
                    </ExpandMore>
                    </div>
                    <Collapse in={expanded} timeout="auto" unmountOnExit>
                        <CardContent style={{maxHeight:200, overflow:"auto"}}>
                                <Messages messages={prop.log} />
                        </CardContent>
                    </Collapse>
                    <div/>
                    <div style={{display:"flex", flexDirection:"row-reverse", alignItems:"center", justifyContent:"space-between", paddingTop:"24px"}}>
                        <Button variant="outlined" disabled={!prop.finished} onClick={()=> {
                            setProgress(-1)
                            prop.setOpenPopUps(false)
                            navigate('.', { replace: true })
                        }}>OK</Button>
                    </div>
                </CardContent>
            </Card>
        </Dialog>
    )
}

const Messages = ({ messages }: { messages:string[] }) => {
    const messagesEndRef = useRef<HTMLInputElement>(null);
    const scrollToBottom = () => {
        messagesEndRef?.current?.scrollIntoView();
    };
    useEffect(scrollToBottom, [messages]);

    return (
        <div className="messagesWrapper">
            {messages.map((message,index) => (
                <Typography key={index}>{message}</Typography>
            ))}
            <div ref={messagesEndRef} />
        </div>
    );
};

function extractProgress(str: string){
    if(str.startsWith("Installation Progress: ")){
        return Number(str.split(": ")[1])*100;
    }
    return -1
}

type ConfigProps = {
    app:string
    openPopUps: boolean
    setOpenPopUps: React.Dispatch<React.SetStateAction<boolean>>
    key:string
    onInstall:(removePlugIns:boolean)=>undefined
}
export function InstallationConfigPopUp(prop: ConfigProps){
    const [removePlugIn, setRemovePlugIn] = useState(false)
    return (<Dialog open={prop.openPopUps} fullScreen={false}>
                <Card style={{width:"300px"}}>
                    <CardContent>
                        <Typography variant="h5" component="div">
                            {prop.app}
                        </Typography>
                        <FormControlLabel control={<Checkbox onChange={()=>setRemovePlugIn(!removePlugIn)} />} label="Remove Plug-ins" />
                        <div style={{
                            display: "flex",
                            flexDirection: "row",
                            alignItems: "center",
                            justifyContent: "flex-end",
                            paddingTop: "24px",
                            gap:"12px"
                        }}><Button variant="outlined" onClick={()=> {
                            prop.setOpenPopUps(false)
                            setRemovePlugIn(false)
                        }}>Cancel</Button>
                            <Button variant="outlined" onClick={()=> {
                                prop.setOpenPopUps(false)
                                prop.onInstall(removePlugIn)
                                setRemovePlugIn(false)
                            }}>Install</Button>
                        </div>
                    </CardContent>
                </Card>
            </Dialog>
    );
}