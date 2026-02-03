import { useState } from "react";
import { MainPage } from "./components/MainPage";
import { Game } from "./components/Game";

type Screen = "main" | "game";

function App() {
  const [currentScreen, setCurrentScreen] = useState<Screen>("main");

  const handleStartGame = () => {
    setCurrentScreen("game");
  };

  const handleExitGame = () => {
    setCurrentScreen("main");
  };

  if (currentScreen === "game") {
    return <Game onExit={handleExitGame} />;
  }

  return <MainPage onStartGame={handleStartGame} />;
}

export default App;
