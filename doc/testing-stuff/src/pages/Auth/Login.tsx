import { LockOutlined } from "@mui/icons-material";
import {
    Container,
    CssBaseline,
    Box,
} from "@mui/material";
import { useState } from "react";
import { Link } from "react-router-dom";

const Login = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    
    const handleLogin = async () => {};

    return (
        <>
            <Container maxWidth="xs">
                <CssBaseline />
                <Box
                    sx={{
                        mt:20,
                        display: 
                    }}
            </Container>
        </>
    );

};

export default Login;