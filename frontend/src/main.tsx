import React from 'react'
import ReactDOM from 'react-dom/client'
import {
    createBrowserRouter, LoaderFunction,
    RouterProvider,
} from "react-router-dom";
import './index.css'
import {AppListLoader, MyApps} from "./routes/myApps.tsx";
import {Homepage} from "./routes/homepage.tsx";
import {Devices, DevicesLoader} from "./routes/devices.tsx";
import {Register} from "./routes/register.tsx";
import {Enrolled, UDIDLoader} from "./routes/enrolled.tsx";


const router = createBrowserRouter([
    {
        path: "",
        element: <Homepage/>,
    },
    {
        path: "register",
        element: <Register />,
    },
    {
        path: "devices",
        element: <Devices/>,
        loader: DevicesLoader as unknown  as LoaderFunction,

    },
    {
        path: "devices/:udid",
        element: <MyApps/>,
        loader: AppListLoader as unknown  as LoaderFunction,
    },
    {
        path: "enrolled/:enrolled",
        element: <Enrolled/>,
        loader: UDIDLoader as unknown  as LoaderFunction,
    },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
      <RouterProvider router={router} />
  </React.StrictMode>,
)
