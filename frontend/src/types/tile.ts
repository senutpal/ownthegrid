export interface Tile {
  id: number;
  x: number;
  y: number;
  ownerId: string | null;
  ownerUsername: string | null;
  ownerColor: string | null;
  claimedAt: string | null;
}

export type TileMap = Map<number, Tile>;
