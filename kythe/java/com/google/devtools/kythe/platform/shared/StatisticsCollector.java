/*
 * Copyright 2014 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.google.devtools.kythe.platform.shared;

/**
 * Allows different analysis drivers to plug-in different statistics collectors
 * to the {@link com.google.devtools.kythe.platform.java.JavacAnalyzer JavacAnalyzer}.
 */
public interface StatisticsCollector {
  /**
   * Increments the named counter by one.
   */
  void incrementCounter(String name);

  /**
   * Increments the named counter by specified amount.
   */
  void incrementCounter(String name, int amount);
}
