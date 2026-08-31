export type MatchStatus = 'scheduled' | 'live' | 'finished' | 'postponed';

export interface Club {
  id: number;
  name: string;
  short_name: string;
  country: string;
  logo_url: string;
  stadium: string;
  founded: number;
}

export interface Match {
  id: number;
  provider: string;
  external_id: string;
  competition_id: number;
  season: string;
  home_club_id: number;
  away_club_id: number;
  match_date: string;
  status: MatchStatus;
  home_score: number;
  away_score: number;
  venue: string;
  minute: number;
  home_club?: Club;
  away_club?: Club;
  competition?: { id: number; name: string };
}

export interface MatchEvent {
  id: number;
  match_id: number;
  minute: number;
  event_type: string;
  team_id: number;
  player_name: string;
  detail: string;
  assist_player_name?: string;
}

export interface MatchLineup {
  id: number;
  match_id: number;
  team_id: number;
  player_id: number;
  player_name: string;
  position: string;
  jersey_number: number;
  is_starter: boolean;
}

export interface MatchStat {
  id: number;
  match_id: number;
  team_id: number;
  stat_type: string;
  value: string;
}

export interface Player {
  id: number;
  name: string;
  position: string;
  jersey_number: number;
  nationality: string;
  photo_url: string;
  rating: number;
  club_id: number;
}

export interface Standing {
  position: number;
  club_id: number;
  club_name: string;
  logo_url: string;
  played: number;
  won: number;
  drawn: number;
  lost: number;
  goals_for: number;
  goals_against: number;
  goal_difference: number;
  points: number;
  form: string;
}

export interface TopScorer {
  player_id: number;
  player_name: string;
  club_name: string;
  goals: number;
  assists: number;
}

export interface User {
  id: number;
  email: string;
  username: string;
  role: string;
}

export interface Paginated<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
}
