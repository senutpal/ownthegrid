import { memo } from 'react';

import type { Tile as TileType } from '../../types/tile';
import styles from './Tile.module.css';

interface TileProps {
  tile: TileType;
  isMyTile: boolean;
  onClick: (id: number) => void;
  disabled: boolean;
}

export const Tile = memo<TileProps>(({ tile, isMyTile, onClick, disabled }) => {
  const isClaimed = Boolean(tile.ownerId);

  return (
    <div
      className={[
        styles.tile,
        isClaimed ? styles.claimed : styles.unclaimed,
        isMyTile ? styles.mine : '',
        disabled ? styles.disabled : '',
      ].join(' ')}
      style={isClaimed && tile.ownerColor ? { backgroundColor: tile.ownerColor } : undefined}
      onClick={() => !disabled && onClick(tile.id)}
      data-tile-id={tile.id}
      title={tile.ownerUsername ?? 'Unclaimed'}
      role="button"
      aria-label={`Tile ${tile.x},${tile.y} — ${tile.ownerUsername ?? 'unclaimed'}`}
    />
  );
});

Tile.displayName = 'Tile';
