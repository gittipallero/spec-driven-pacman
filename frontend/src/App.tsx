import { MainPage } from "./components/MainPage";

function App() {
  const handleStartGame = () => {
    // TODO: Navigate to game screen when game feature is implemented
    console.log("Starting game...");
  };

  return <MainPage onStartGame={handleStartGame} />;
}

export default App;
