// index.ts
import axios from 'axios';
let accessToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNYXBDbGFpbXMiOnsiYXVkIjoiaHR0cDovLzIwLjI0NC41Ni4xNDQvZXZhbHVhdGlvbi1zZXJ2aWNlIiwiZW1haWwiOiJwYXJ2YWdnYXJ3YWwxMzBAZ21haWwuY29tIiwiZXhwIjoxNzc4NzU4OTU1LCJpYXQiOjE3Nzg3NTgwNTUsImlzcyI6IkFmZm9yZCBNZWRpY2FsIFRlY2hub2xvZ2llcyBQcml2YXRlIExpbWl0ZWQiLCJqdGkiOiIyZjQwM2EyNC1lZTExLTQ0NGQtYTEzMi0yYjgyNDcyNjBhM2MiLCJsb2NhbGUiOiJlbi1JTiIsIm5hbWUiOiJwYXJ2IGFnZ2Fyd2FsIiwic3ViIjoiNzViOGI3NjQtMmM2MS00Mzg4LWFhOTYtZDNjMTI2ZWMxODA1In0sImVtYWlsIjoicGFydmFnZ2Fyd2FsMTMwQGdtYWlsLmNvbSIsIm5hbWUiOiJwYXJ2IGFnZ2Fyd2FsIiwicm9sbE5vIjoiMTIzMTA3MDAiLCJhY2Nlc3NDb2RlIjoiVFJ2WldxIiwiY2xpZW50SUQiOiI3NWI4Yjc2NC0yYzYxLTQzODgtYWE5Ni1kM2MxMjZlYzE4MDUiLCJjbGllbnRTZWNyZXQiOiJId3hSQVRlaFlydmpiVUROIn0.BjZEnmFmkOtpw5A_uym6YpYwCGO6DXWizqZ6ceYmWSM';
export const initLogger = (token) => {
    accessToken = token;
};
export async function Log(stack, level, pkg, message) {
    if (!accessToken) {
        console.warn("Logger warning: Access token not set.");
        return;
    }
    try {
        await axios.post('http://4.224.186.213/evaluation-service/logs', {
            stack: stack,
            level: level,
            package: pkg,
            message: message
        }, {
            headers: {
                'Authorization': `Bearer ${accessToken}`,
                'Content-Type': 'application/json'
            }
        });
    }
    catch (error) {
        console.error("Logging Middleware Error:", error.response?.data || error.message);
    }
}
