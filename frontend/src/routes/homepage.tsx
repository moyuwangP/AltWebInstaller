import {useNavigate} from 'react-router-dom';
import {useEffect} from "react";
import {UDID, UDIDLoader} from "./enrolled.tsx";


export function Homepage(){
    const navigate = useNavigate();
    const udidString = localStorage.getItem("udid")

    useEffect(() => {
        if(isIOS()){
            if(udidString != null){
                navigate(`/devices/${udidString}`)
            }else{
                UDIDLoader().then(function (udid:UDID){
                    if(udid.ok) {
                        localStorage.setItem("udid", udid.udid)
                        navigate(`/devices/${udid.udid}`)
                    }else{
                        navigate("/register");
                    }
                })
            }

        }else{
            navigate("/devices")
        }
    }, [navigate]);

    return (
        <div></div>
    )
}

function isIOS() {

    if (/iPad|iPhone|iPod/.test(navigator.userAgent)) {
        return true;
    } else {
        return navigator.maxTouchPoints > 1 &&
            /Macintosh/.test(navigator.userAgent);
    }
}