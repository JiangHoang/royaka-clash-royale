import { Navigate } from "react-router-dom";
import { hasStoredSession } from "../utils/session";

const PrivateRoute = ({ children }) => {
    if (!hasStoredSession()) {
        return <Navigate to="/auth" replace />;
    }

    return children;
};

export default PrivateRoute;
