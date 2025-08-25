import { useNavigate } from "react-router";
import { useLogout } from "../auth/hooks";

const NavigationBar = (): React.JSX.Element => {

  const { logout } = useLogout();

  const navigate = useNavigate();

  const handleDSAClick = () => {
    // Navigate to main page (implementation to be added)
    navigate("/dashboard");
  };

  const handleLogout = () => {
    // Logout process (implementation to be added)
    logout();
    navigate("/login");
  };

  return (
    <div className="bg-blue-500 text-white px-6 py-4 flex justify-between items-center">
      <button
        key="dsa-button"
        onClick={handleDSAClick}
        className="text-2xl font-bold hover:opacity-80 transition-opacity"
      >
        DSA
      </button>
      <button
        key="logout-button"
        onClick={handleLogout}
        className="hover:bg-blue-600 px-4 py-2  rounded transition-colors"
      >
        Logout
      </button>
    </div>
  )
}

export default NavigationBar;
