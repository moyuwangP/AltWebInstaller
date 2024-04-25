import {useLoaderData, useNavigate} from "react-router-dom";
import Typography from "@mui/material/Typography";

type device ={
    udid:string
    mac_address:string
    online:boolean
}
export async function DevicesLoader()  {
    return fetch(`/api/devices`)
        .then(function (response) {
            if (!response.ok) {
                throw Error("SERVERR");
            }
            return response.json();
        }).then(function (nodes: device[]) {
            return nodes
        }).catch(function (error) {
            throw Error(error.toString());
        })
}

export function Devices(){
    const devices  = useLoaderData() as device[];
    const navigate = useNavigate()

    const nodeList = devices
        .map((node: device) => (
            <div style={{
                borderRadius: 30,
                height: "7rem",
                backgroundColor: node.online?
                    "rgba(15,126,130,0.35)":
                    "rgb(205,205,205)",
                marginTop:16,
                display: 'flex',
                alignItems:'center',
                paddingLeft:16,
                paddingRight:16,
                justifyContent:'space-between',
                cursor: "pointer",
            }}
                 onClick={() => navigate(node.udid)}
                 key={node.udid}
            >
                <div style={{display: 'flex',flexDirection:"column"}}>
                    <Typography fontFamily={"CircularSpBold"}>{node.udid}</Typography>
                    <Typography fontFamily={"CircularSpBold"}>{node.mac_address}</Typography>
                </div>
            </div>
        ))
    return (
        <div style={{
            padding: "10px",
            maxWidth: "430px ",
            marginLeft: "auto",
            marginRight: "auto",
            width: "100%",
            height: "100%",
            boxSizing: "border-box"
        }}>
            <Typography fontSize={'2.5rem'} fontFamily={"CircularSpBold"}>Devices</Typography>
            <Typography fontSize={'1.5rem'} fontFamily={"CircularSpBold"}>Paired</Typography>
            {nodeList}
        </div>
    )
}