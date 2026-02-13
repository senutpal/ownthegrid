import { useMemo } from 'react';

import { useBoardStore } from '../../store/boardStore';
import { useUserStore } from '../../store/userStore';
import { useBoard } from '../../hooks/useBoard';
import { Tile } from './Tile';
import styles from './Grid.module.css';

export const Grid = () => {
  const tiles = useBoardStore((state) => state.tiles);
  const isLoaded = useBoardStore((state) => state.isLoaded);
  const gridWidth = useBoardStore((state) => state.gridWidth);
  const gridHeight = useBoardStore((state) => state.gridHeight);
  const { currentUser } = useUserStore();
  const { handleClaim } = useBoard();

  const tileElements = useMemo(() => {
    if (!isLoaded) return null;
    return Array.from(tiles.values()).map((tile) => (
      <Tile
        key={tile.id}
        tile={tile}
        isMyTile={tile.ownerId === currentUser?.id}
        onClick={handleClaim}
        disabled={false}
      />
    ));
  }, [tiles, currentUser?.id, handleClaim, isLoaded]);

  if (!isLoaded) {
    return <div className={styles.loading}>Loading board...</div>;
  }

  // Calculate aspect ratio and max-width dynamically based on grid dimensions
  const aspectRatio = `${gridWidth} / ${gridHeight}`;
  const maxWidth = `min(92vw, calc(90vh * ${gridWidth} / ${gridHeight}))`;

  return (
    <div
      className={styles.grid}
      style={{
        gridTemplateColumns: `repeat(${gridWidth}, 1fr)`,
        gridTemplateRows: `repeat(${gridHeight}, 1fr)`,
        aspectRatio,
        maxWidth,
      }}
    >
      {tileElements}
    </div>
  );
};
