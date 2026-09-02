
/** @param {{ message?: string }} props */
export function LoadingScreen({ message = "Loading..." }) {
    return (
        <div className="min-h-screen bg-gradient-to-b from-blue-900 via-blue-700 to-blue-900 flex items-center justify-center">
            <div className="rounded-xl border-4 border-yellow-400 bg-blue-800 px-8 py-5 text-xl text-white shadow-xl">
                {message}
            </div>
        </div>
    );
}
