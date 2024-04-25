import Typography from "@mui/material/Typography";
import {useLoaderData} from "react-router-dom";

export type UDID = {
    udid:string
    ok: boolean
    msg: string
}
export async function UDIDLoader()  {
    return fetch(`/api/udid`)
        .then(function (response) {
            if (!response.ok) {
                throw Error("SERVERR");
            }
            return response.json();
        }).then(function (nodes: UDID) {
            nodes.ok = true
            nodes.msg = ""
            return nodes
        }).catch(function (error) {
            return {
                udid:"",
                ok: false,
                msg: error.toString()
            } as UDID
        })
}

export function Enrolled(){
    const  udid  = useLoaderData() as UDID;
    console.log(udid)
    return (
        <div style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            flexDirection:"column",
            width: "100vw",
            height: "80vh"
        }}>
            <Typography fontSize={'2.5rem'} fontFamily={"CircularSpBold"}>You Are All Set</Typography>
            <Typography fontSize={'1.5rem'} fontFamily={"CircularSpBold"}>Go Back to AltWebInstaller</Typography>
            <Typography fontSize={'1.5rem'} fontFamily={"CircularSpBold"}>to Manage Your App</Typography>
            <Typography fontSize={'.5rem'} fontFamily={"CircularSpBold"}>Your UDID:{udid.udid}</Typography>
        </div>
    )
}