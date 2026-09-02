import { Routes, Route } from "react-router-dom";
import Entry from "@/pages/entry";
import GameSimple from "@/pages/game-simple";
import GameEnhanced from "@/pages/game-enhanced";
import Lobby from "@/pages/lobby";
import Auth from "@/pages/auth";
import CardDesk from "@/pages/card-desk";
import PrivateRoute from "@/app/router/PrivateRoute";
import { WebSocketProvider } from "@/app/providers/websocket";

function App() {
    return (
        <WebSocketProvider>
            <Routes>
                <Route path="/" element={<Entry />} />
                <Route path="/index.html" element={<Entry />} />
                <Route path="/auth" element={<Auth />} />
                <Route path="/lobby" element={<PrivateRoute><Lobby /></PrivateRoute>} />
                <Route path="/game-simple" element={<PrivateRoute><GameSimple /></PrivateRoute>} />
                <Route path="/game-enhanced" element={<PrivateRoute><GameEnhanced /></PrivateRoute>} />
                <Route path="/card-desk" element={<PrivateRoute><CardDesk /></PrivateRoute>} />
            </Routes>
        </WebSocketProvider>
    );
}

export default App;
