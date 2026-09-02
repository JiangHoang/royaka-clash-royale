import { Navigate } from "react-router-dom";
import { hasStoredSession } from "@/shared/lib";

const PrivateRoute = ({ children }) => {
    if (!hasStoredSession()) {
        return <Navigate to="/auth" replace />;
    }

    return children;
};

export default PrivateRoute;
