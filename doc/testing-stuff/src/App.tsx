import  React, { useState, useEffect } from "react";
import { Route, Switch } from "react-router-dom";
import { useAppDispatch } from '@hooks';
import * as actions from 'actions';

import { updateAPIConfig } from './api/base'
import { isValidToken } from '';

import Layout from 'containers/Layout';

import Landing from "./pages/Landing";
import Navbar from "./components"
import Home from "./pages.Home";
import Login from "./pages/Auth/Login";
import Register from "./pages/Auth/Register";

function App() {
    return (
        <>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
            </Routes>
        </>
    );
}

export default App;